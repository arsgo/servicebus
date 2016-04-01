/*
提供服务提供商功能
1. 下载服务提供商配置文件，并根据服务模式分组，检查当前是否是独占服务器，如果是检查当前服务是否正确，正确则退出，否则轮循每个配置
2. 检查当前机器IP是否与配置相符，不符则重复执行步骤2检查下一个配置，否则转到步骤3
3. 创建当前服务配置节点
4. 检查当前服务的数量是否超过配置数据，如果超过则删除当前节点
5. 检查当前配置是否是独占，如果是则返回状态，不再继续绑定服务，否则转到步骤2绑定下一服务
6. 标记当前服务器是独占还是共享，如果是共享则转到步骤7执行，否则转到步骤8执行
7. 监控所有独占服务变化，变化后，重新绑定当前服务，绑定成功后删除所有共享服务
8. 监控服务配置信息变化，变化后执行步骤1


*/

package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

func (d *serviceProvider) WatchServiceConfigChange() {
	path := d.dataMap.Translate(serviceConfig)
START:
	if !zkClient.ZkCli.Exists(path) {
		time.Sleep(time.Second * 30)
        d.Log.Infof("sp config not exists:%s",path)
		goto START
	}
	d.rebindAllServices()

	changes := make(chan string, 10)
	go zkClient.ZkCli.WatchValue(path, changes)

	for {
		select {
		case <-changes:
			{
				d.Log.Info("recv service provider config changed notices")
				d.rebindAllServices()
			}
		}
	}

}

func (d *serviceProvider) rebindAllServices() {
	d.lk.Lock()
	defer d.lk.Unlock()
	aloneService, sharedService := d.downloadServiceProvicerConfig()
	if strings.EqualFold(d.mode, eModeAlone) {
		goon, err := d.checkAloneService(d.services.services, aloneService)
		if !goon {
			return
		}
		if err != nil {
			d.Log.Error(err)
		}
	}
	config := d.bindServices(aloneService)
	if len(config.services) > 0 {
		d.deleteSharedSevices(config.services)
		d.mode = eModeAlone
		d.services = config
		return
	}

	config = d.bindServices(sharedService)
	d.mode = eModeShared
	d.deleteSharedSevices(config.services)
	d.services = config
}

//UpdateServiceData  update service consumer data
func (d *serviceProvider) UpdateServiceData(serviceName string, ndata utility.DataMap) {
	nmap := d.dataMap.Merge(ndata)
	nmap.Set("serviceName", serviceName)
	nmap.Set("last", fmt.Sprintf("%d", time.Now().Unix()))
	path := nmap.Translate(serviceProviderPath)
	zkClient.ZkCli.UpdateValue(path, nmap.Translate(serviceProviderValue))
}

//ClearServiceProviders clear all service providers
func (d *registerCenter) ClearServiceProviders() (string, error) {
	path := d.dataMap.Translate(serviceRoot)
	err := zkClient.ZkCli.Delete(path)
	if err != nil {
		return "", err
	}
	zkClient.ZkCli.CreatePath(path, "")
	return path, nil
}

func init() {
	var err error
	ServiceProvider = &serviceProvider{}
	ServiceProvider.dataMap = utility.NewDataMap()
	ServiceProvider.dataMap.Set("domain", zkClient.Domain)
	ServiceProvider.dataMap.Set("ip", zkClient.LocalIP)
	ServiceProvider.Log, err = logger.New("service provider", true)
	ServiceProvider.services = &providerServicesConfig{}
	ServiceProvider.services.services = make(map[string]*providerService, 0)
	if err != nil {
		log.Println(err)
	}
}

//-----------------------------------------------------------------------------------------
func (d *serviceProvider) deleteSharedSevices(svs map[string]*providerService) {
	for i := range d.services.services {
		if _, ok := svs[i]; ok {
			continue
		}
		nmap := d.dataMap.Copy()
		nmap.Set("serviceName", i)
		path := nmap.Translate(serviceProviderPath)
		zkClient.ZkCli.Delete(path)
		d.Log.Infof("delete service:%s", i)
	}

}

func (d *serviceProvider) downloadServiceProvicerConfig() (aloneService map[string]*providerService,
	sharedService map[string]*providerService) {
	aloneService = make(map[string]*providerService)
	sharedService = make(map[string]*providerService)
	svs := []*providerService{}
	path := d.dataMap.Translate(serviceConfig)
	values, err := zkClient.ZkCli.GetValue(path)
	if err != nil {
		d.Log.Errorf("error:%s,%+v", path, err)
		return
	}
	err = json.Unmarshal([]byte(values), &svs)
	if err != nil {
		d.Log.Infof("services provider config unmarshal error:%+v", err)
	}
	for _, v := range svs {
		if strings.EqualFold(v.Mode, eModeAlone) {
			aloneService[v.Name] = v
		} else if strings.EqualFold(v.Mode, eModeShared) {
			sharedService[v.Name] = v
		}
	}
	return
}
func (d *serviceProvider) checkAloneService(current map[string]*providerService, configs map[string]*providerService) (ct bool, err error) {
	ct = true
	err = nil
	if len(current) > 1 {
		for i := range current {
			nmap := d.dataMap.Copy()
			nmap.Set("serviceName", i)
			path := nmap.Translate(serviceProviderPath)
			zkClient.ZkCli.Delete(path)
		}
		return
	} else if len(current) < 1 {
		return
	}
	for i, v := range current {
		if cc, ok := configs[i]; ok && (strings.EqualFold(v.getUNIQ(), cc.getUNIQ()) || checkIP(configs[i].IP)) {
			nmap := d.dataMap.Copy()
			nmap.Set("serviceName", i)
			path := nmap.Translate(serviceProviderPath)
			if zkClient.ZkCli.Exists(path) {
				ct = false
				return
			}
		}
	}
	return
}

func checkIP(origin string) bool {
	ips := fmt.Sprintf(",%s,", origin)
	llocal := fmt.Sprintf(",%s,", zkClient.LocalIP)
	return strings.Contains(ips, llocal)
}

func (d *serviceProvider) bindServices(services map[string]*providerService) (psconfig *providerServicesConfig) {
	psconfig = &providerServicesConfig{}
	psconfig.services = make(map[string]*providerService, 0)
	for sv, config := range services {
		var (
			err error
		)
		if v, ok := d.services.services[sv]; !ok {
			err = d.bindService(sv, config)
		} else if !strings.EqualFold(config.getUNIQ(), v.getUNIQ()) {
			err = d.bindService(sv, config)
		}
		if err == nil {
			psconfig.services[sv] = config
		}
		if strings.EqualFold(config.Mode, eModeAlone) && err == nil {
			return
		}
	}
	return
}

//BindServices  绑定服务
func (d *serviceProvider) bindService(serviceName string, config *providerService) (err error) {
	nmap := d.dataMap.Copy()
	nmap.Set("serviceName", serviceName)
	nmap.Set("last", fmt.Sprintf("%d", time.Now().Unix()))
	path := nmap.Translate(serviceProviderPath)
	if !checkIP(config.IP) {
		if zkClient.ZkCli.Exists(path) {
			err = zkClient.ZkCli.Delete(path)
		}
		return errors.New("ip not match")
	}
	if zkClient.ZkCli.Exists(path) {
		return
	}

	_, err = zkClient.ZkCli.CreateTmpNode(path, nmap.Translate(serviceProviderValue))
	if err != nil {
		return
	}
	d.Log.Infof("::register service:%s", serviceName)
	return
}
