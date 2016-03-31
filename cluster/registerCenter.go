/*
 提供注册中心与ZK之前的管理包括:
 1. 将当前服务器IP等信息注册到ZK
 2. 检查当前服务器是master还是salve
 3. 如果是savle则：
         a.监控注册中心服务器变化，变化后下载服务器列表, 并转到第2步执行
         b.下载服务提供商列表到本地，并监控服务列表变化，变化后重复执行b，
    如果是master则转到第4步执行
 4. 更新ZK中的服务节点值为master
 5. 下载当前域下所有服务提供商地址
 6. 生成服务提供商列表，并发布到服务列表地址
 7. 更新本地服务列表
 8. 监控服务提供商变化，变化后转到第5步执行

*/

package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
    "github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

//Bind  bind RC Server to zookeeper cluster
func (d *registerCenter) Bind(log *logger.Logger) (path string, err error) {
	d.log=log
    err = d.createPath()
	if err != nil {
		d.log.Error(err)
		return
	}
	go d.watchMasterChange()
	go d.watchCurrentPath()
	return
}

//WatchJobChange 监控JOB配置变化
func (d *registerCenter) WatchJobChange(callback func(config *JobConfigs, err error)) error {
	d.jobCallback = callback
	if !d.IsMasterServer {
		return nil
	}
	if callback == nil {
		return errors.New("callback func must value")
	}
	configPath := d.dataMap.Translate(jobConfigPath)
	jobChanges := make(chan string, 10)
	d.log.Info("::start watch job config changes")
	go func() {
		var waitNotify bool
	START:
		if !zkClient.ZkCli.Exists(configPath) {
			waitNotify = true
			time.Sleep(time.Second * 10)
			goto START
		}
		if waitNotify {
			d.log.Info("recv job config changed notices")
		}
		callback(d.getJobConfigs())
		go zkClient.ZkCli.WatchValue(configPath, jobChanges)
		for {
			select {
			case <-jobChanges:
				{
					d.log.Info("recv job config changed")
					configs, err := d.getJobConfigs()
					if err != nil {
						d.log.Infof("error:%+v",err)
					}
                    callback(configs, err)
				}
			}

		}
	}()
	return nil
}

//WatchServiceListChange 监控服务列表变化
func (d *registerCenter) WatchServicesChange(callback func(services map[string][]string, err error)) error {
	if callback == nil {
		return errors.New("callback func must value")
	}
	d.log.Info("::start watch services list changes")
	rootPath := d.dataMap.Translate(servicePublishPath)
	changes := make(chan string, 10)
	go func() {
		var waitNotify bool
	START:
		if !zkClient.ZkCli.Exists(rootPath) {
			waitNotify = true
			time.Sleep(time.Second * 10)
			goto START
		}
		if waitNotify {
			d.log.Info("recv service changed notices")
		}

		d.setServiceUpdateTime()
		callback(d.downloadProviders())
		go zkClient.ZkCli.WatchValue(rootPath, changes)
		for {
			select {
			case <-changes:
				{
					d.log.Info("recv service changed notices")
					d.setServiceUpdateTime()
					callback(d.downloadProviders())
				}
			}
		}
	}()
	return nil
}

//NewRegisterCenter create RC
func NewRegisterCenter(domain string, ip string) *registerCenter {
	rc := &registerCenter{}
	rc.dataMap = utility.NewDataMap()
	rc.dataMap.Set("domain", zkClient.Domain)
	rc.dataMap.Set("ip", zkClient.LocalIP)
	rc.dataMap.Set("type", "salve")
	return rc
}

func init() {
	RegisterCenter = NewRegisterCenter(zkClient.Domain, zkClient.LocalIP)
}

//------------------------------------------------------------------------------------------------

//GetConfigs Get job Config
func (d *registerCenter) getJobConfigs() (defConfigs *JobConfigs, err error) {
	defConfigs = &JobConfigs{}
	defConfigs.Jobs = make(map[string]JobConfigItem)
	configPath := d.dataMap.Translate(jobConfigPath)
	if !zkClient.ZkCli.Exists(configPath) {
        d.log.Info("job config path not exists")
		return
	}
	value, err := zkClient.ZkCli.GetValue(configPath)
	if err != nil {
		d.log.Infof("get job config error:%+v", err)
		return
	}
	err = json.Unmarshal([]byte(value), &defConfigs.Jobs)
	if err != nil {
		d.log.Infof("job config unmarshal error:%+v", err)
	}
	return
}

//WatchMasterChange  watch whether the master server is changed
func (d *registerCenter) watchMasterChange() {

	if d.IsMasterServer {
		go d.watchServiceProviderChange()
		return
	}
	d.log.Info("::start watch master changes")
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
				if m := d.isMaster(d.Path, data); m != d.IsMasterServer && m {
					d.setOnline(d.Path, true)
					d.serviceChange()
					go d.watchServiceProviderChange()
					break EndLoop
				}
			}
		}
	}
}

func (d *registerCenter) watchCurrentPath() {
	tp := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-tp.C:
			if !zkClient.ZkCli.Exists(d.Path) {
				d.log.Info("RC not exists")
				d.createPath()
			}
		}
	}

}

//WatchServiceProviderChange watch whether any service privider is changed
func (d *registerCenter) watchServiceProviderChange() error {
	if !d.IsMasterServer {
		return nil
	}
	d.log.Info("::start watch service providers changes")
	var changes chan []string
	changes = make(chan []string, 10)
	childrenChanged := make(chan []string, 10)
	rootPath := d.dataMap.Translate(serviceRoot)
	if !zkClient.ZkCli.Exists(rootPath) {
		d.log.Infof("services node not exists:%s\r\n", rootPath)
		time.Sleep(time.Second * 10)
		return d.watchServiceProviderChange()
	}
	d.serviceChange()

	go zkClient.ZkCli.WatchChildren(rootPath, changes)
	go func() {
		children, err := zkClient.ZkCli.GetChildren(rootPath)
		if err != nil {
			d.log.Error(err)
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
				d.serviceChange()
			}
		}
	}
}

func (d *registerCenter) publishServices() error {
	if !d.IsMasterServer {
		return errors.New("必须是master才能发布服务列表")
	}
	providers, err := d.downloadProviders()
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

//DownloadServiceProviders download the list of service providers
func (d *registerCenter) downloadProviders() (ServiceProviderList, error) {
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
func (d *registerCenter) createPath() (err error) {
	d.Path, err = zkClient.ZkCli.CreateSeqNode(d.dataMap.Translate(dsbServerPath), d.dataMap.Translate(dsbServerValue))
	if err != nil {
		return
	}
	children, err := zkClient.ZkCli.GetChildren(d.dataMap.Translate(dsbServerRoot))
	if err != nil {
		return
	}
	d.setOnline(d.Path, d.isMaster(d.Path, children))
	d.updateNodeValue()
	return
}

//isMaster  check whether the current server is the master server
func (d *registerCenter) isMaster(path string, servers []string) bool {
	sort.Sort(sort.StringSlice(servers))
	return len(servers) == 0 || strings.HasSuffix(path, servers[0])
}

func (d *registerCenter) serviceChange() {
	d.log.Info("service provider has changed")
	d.setServiceUpdateTime()
	d.updateNodeValue()
	err := d.publishServices()
	if err != nil {
		d.log.Infof("service publish failed:%s", err.Error())
	} else {
		d.log.Info("service publish success")
	}
	time.Sleep(time.Second)
}

//UpdateMasterValue set current as master
func (d *registerCenter) updateNodeValue() {
	pvalue := d.dataMap.Translate(dsbServerValue)
	err := zkClient.ZkCli.UpdateValue(d.Path, pvalue)
	if err != nil {
		fmt.Printf("更新NODE发生错误:%s,%s", d.Path, err.Error())
	}
}

//GetMasterValue get mask value
func (d *registerCenter) getValue() (*dsbCenterNodeValue, string) {
	value, _ := zkClient.ZkCli.GetValue(d.Path)
	values := []byte(value)
	dsb := &dsbCenterNodeValue{}
	err := json.Unmarshal(values, dsb)
	if err != nil {
		d.log.Infof("get master value error,[%s],[%s]", value, d.Path)
		d.log.Error(err.Error())
	}
	return dsb, value
}
func (d *registerCenter) setServiceUpdateTime() {
	d.LastPublish = time.Now().Unix()
	d.dataMap.Set("last", fmt.Sprintf("%d", d.LastPublish))
	d.dataMap.Set("pst", fmt.Sprintf("%d", d.LastPublish))
}
func (d *registerCenter) setOnline(path string, master bool) {
	d.Path = path
	d.IsMasterServer = master
	d.OnlineTime = time.Now().Unix()
	d.dataMap.Set("path", d.Path)
	d.dataMap.Set("online", fmt.Sprintf("%d", d.OnlineTime))
	d.dataMap.Set("last", fmt.Sprintf("%d", d.OnlineTime))
	d.log.Infof("online:%s", path)
	if d.IsMasterServer {
		d.log.Info("server is master")
		d.dataMap.Set("type", "master")
		d.Server = "master"
		d.WatchJobChange(d.jobCallback)
	} else {
		d.log.Info("server is salve")
		d.dataMap.Set("type", "salve")
		d.Server = "salve"
	}
}
