package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/colinyl/lib4go/utility"
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
	d.IsMasterServer = d.isMaster(children)
	if d.IsMasterServer {
		d.UpdateMasterValue()
	}
	go d.WatchMasterChange()
	return d.Path, nil
}

//isMaster  check whether the current server is the master server
func (d *registerCenter) isMaster(servers []string) bool {
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

//PublishServiceList To update the list of services
func (d *registerCenter) PublishServiceList() error {
	if !d.IsMasterServer {
		return errors.New("必须是master才能发布服务列表")
	}
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

//WatchServiceListChange 监控服务列表变化
func (d *registerCenter) WatchServiceListChange(callback func(services map[string][]string, err error)) error {
	changes := make(chan []string, 10)
	rootPath := d.dataMap.Translate(servicePublishPath)
	go func() {
		if callback == nil {
			return
		}
		callback(d.DownloadServiceProviders())
		go zkClient.ZkCli.WatchChildren(rootPath, changes)
		for {
			select {
			case <-changes:
				{
					callback(d.DownloadServiceProviders())
				}
			}
		}
	}()
	return nil
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
				d.PublishServiceList() //更新服务
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
				if m := d.isMaster(data); m != d.IsMasterServer && m {
					d.IsMasterServer = true
					d.UpdateMasterValue()
					d.PublishServiceList() //更新服务
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

//GetMasterValue get mask value
func (d *registerCenter) GetMasterValue() (*dsbCenterNodeValue, string) {
	value, _ := zkClient.ZkCli.GetValue(d.Path)
	values := []byte(value)
	dsb := &dsbCenterNodeValue{}
	err := json.Unmarshal(values, dsb)
	if err != nil {
		Log.Infof("get master value error,[%s],[%s]", value, d.Path)
		Log.Error(err)
	}
	return dsb, value
}
func (d *registerCenter) GetSnap() string {
	data := make(map[string]string)
	data["isMaster"] = fmt.Sprintf("%s", d.IsMasterServer)
	data["path"] = d.Path
	data["last"] =  fmt.Sprintf("%d", d.Last)
	buffer, _ := json.Marshal(data)
	return string(buffer)
}

func NewRegisterCenter(domain string, ip string) *registerCenter {
	rc := &registerCenter{}
	rc.dataMap = utility.NewDataMap()
	rc.dataMap.Set("domain", zkClient.Domain)
	rc.dataMap.Set("ip", zkClient.LocalIP)
	rc.dataMap.Set("isMaster", "false")
	return rc
}

func init() {
	RegisterCenter = NewRegisterCenter(zkClient.Domain, zkClient.LocalIP)
}
