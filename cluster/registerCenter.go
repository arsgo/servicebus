package cluster

import (
	"encoding/json"
	"fmt"
	"github.com/lib4go/lib4go/utility"
	"sort"
	"strings"
	"time"
)

//Bind  bind dsbcenter server to zk
func (d *registerCenter) Bind() (string, error) {
	var err error
	d.Path, err = zkClient.ZkCli.CreateSeqNode(d.dataMap.Translate(dsbServerPath), d.dataMap.Translate(dsbServerValue))
    if err != nil {
		Log.Error(err)
		return "", err
	}

	children, er := zkClient.ZkCli.GetChildren(d.dataMap.Translate(dsbServerRoot))
	if er != nil {
		return "", er
	}
	d.IsMasterServer = d.IsMaster(children)
	if d.IsMasterServer {
		d.UpdateMasterValue()
	}
	go d.WatchMasterChange()
	return d.Path, nil
}

//IsMaster  check whether the current server is the master server
func (d *registerCenter) IsMaster(servers []string) bool {
	sort.Sort(sort.StringSlice(servers))
	return strings.HasSuffix(d.Path, servers[0])
}


//DownloadServiceProviders download the list of service providers
func (d *registerCenter) DownloadServiceProviders() (ServiceProviderList, error) {
	var spList ServiceProviderList = make(map[string][]string)
	serviceList, err := zkClient.ZkCli.GetChildren(d.dataMap.Translate(serviceRoot))
	if err != nil {
		return spList, err
	}

	for _, v := range serviceList {
		nmap := d.dataMap.Copy()
		nmap.Set("serviceName", v)
		providerList, er := zkClient.ZkCli.GetChildren(nmap.Translate(serviceProviderRoot))
		if er != nil {
			return spList, er
		}
		for _, l := range providerList {
			spList.Add(v, l)
		}
	}
	return spList, nil
}

//UpdateServiceList To update the list of services
func (d *registerCenter) UpdateServiceList() error {
	providers, err := d.DownloadServiceProviders()
	if err != nil {
		return err
	}

	data, err := json.Marshal(providers)
	if err != nil {
		return err
	}

	serverListData := string(data)
	path := d.dataMap.Translate(servicePublishPath)
	if zkClient.ZkCli.Exists(path) {
		return zkClient.ZkCli.UpdateValue(path, serverListData)
	}
	return zkClient.ZkCli.CreatePath(path, serverListData)
}

//WatchServiceProviderChange watch whether any service privider is changed
func (d *registerCenter) WatchServiceProviderChange() error {
	if !d.IsMasterServer {
		return nil
	}
	var changes chan []string
	changes = make(chan []string, 10)
	childrenChanged := make(chan []string, 10)
	rootPath := d.dataMap.Translate(serviceRoot)
	go zkClient.ZkCli.WatchChildren(rootPath, changes)
	go func() {
		children, err := zkClient.ZkCli.GetChildren(rootPath)
		if err != nil {
			Log.Error(err)
		}
		for _, v := range children {
			servicePath := fmt.Sprintf("%s/%s/providers", rootPath, v)
			go zkClient.ZkCli.WatchChildren(servicePath, childrenChanged)
		}
	}()
	for {
		select {
		case <-childrenChanged:
			{
				d.UpdateMasterValue()
				d.UpdateServiceList() //更新服务
				time.Sleep(time.Second)
			}
		}
	}

}

//WatchMasterChange  watch whether the master server is changed
func (d *registerCenter) WatchMasterChange() {
	if d.IsMasterServer {
		go d.WatchServiceProviderChange()
		return
	}
	var children chan []string
	children = make(chan []string, 10)
	path := d.dataMap.Translate(dsbServerRoot)
	go zkClient.ZkCli.WatchChildren(path, children)
EndLoop:
	for {
		select {
		case <-children:
			{
				data, _ := zkClient.ZkCli.GetChildren(path)
				if m := d.IsMaster(data); m != d.IsMasterServer && m {
					d.IsMasterServer = true
					d.UpdateMasterValue()
					d.UpdateServiceList() //更新服务
					go d.WatchServiceProviderChange()
					break EndLoop
				}
			}
		}
	}
}

//UpdateMasterValue set current as master
func (d *registerCenter) UpdateMasterValue() {
	if d.IsMasterServer {
		d.dataMap.Set("isMaster", "true")
		d.Last = time.Now().Unix()
		d.dataMap.Set("now", fmt.Sprintf("%d", d.Last))
		zkClient.ZkCli.UpdateValue(d.Path, d.dataMap.Translate(dsbServerValue))
	}
}

//GetMaskerValue get mask value
func (d *registerCenter) GetMaskerValue() (*dsbCenterNodeValue, string) {
	value, _ := zkClient.ZkCli.GetValue(d.Path)
	values := []byte(value)
	dsb := &dsbCenterNodeValue{}
	err := json.Unmarshal(values, dsb)
	if err != nil {
		Log.Infof("get masker value error,[%s],[%s]", value, d.Path)
		Log.Error(err)
	}
	return dsb, value
}

func init() {
	RegisterCenter = &registerCenter{}
	RegisterCenter.dataMap = utility.NewDataMap()
	RegisterCenter.dataMap.Set("domain", zkClient.Domain)
	RegisterCenter.dataMap.Set("ip", zkClient.LocalIP)
	RegisterCenter.dataMap.Set("isMaster", "false")
}
