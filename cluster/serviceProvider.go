package cluster

import (
	"fmt"
	"time"
    "log"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

//BindServices  update consumer data
func (d *serviceProvider) BindServices(serviceNames []string, ndata utility.DataMap) error {
	for _, v := range serviceNames {
		nmap := d.dataMap.Merge(ndata)
		nmap.Set("serviceName", v)
		nmap.Set("now", fmt.Sprintf("%d", time.Now().Unix()))
		path := nmap.Translate(serviceProviderPath)
		_, err := zkClient.ZkCli.CreateTmpNode(path, nmap.Translate(serviceProviderValue))
		if err != nil {
			return err
		}
	}
	return nil
}

//UpdateServiceData  update service consumer data
func (d *serviceProvider) UpdateServiceData(serviceName string, ndata utility.DataMap) {
	nmap := d.dataMap.Merge(ndata)
	nmap.Set("serviceName", serviceName)
	nmap.Set("now", fmt.Sprintf("%d", time.Now().Unix()))
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
     if err!=nil{
        log.Print(err)
    }
}
