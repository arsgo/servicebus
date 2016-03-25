package cluster

import (
	"fmt"
	"github.com/colinyl/lib4go/utility"
	"time"
)

//BindService consumer bind a service
func (d *serviceConsumer) BindService(serviceName string, ndata utility.DataMap) (string, string) {
	nmap := d.dataMap.Merge(ndata)
	nmap.Set("serviceName", serviceName)
	nmap.Set("now", fmt.Sprintf("%d", time.Now().Unix()))
	path := nmap.Translate(serviceConsumerPath)
	if _, ok := d.PathList[path]; !ok {
		zkClient.ZkCli.CreateTmpNode(path, nmap.Translate(serviceConsumerValue))
		d.PathList[path] = path
	}
	return d.PathList[path], nmap.Translate(serviceConsumerValue)
}

//UpdateConsumerData  update consumer data
func (d *serviceConsumer) UpdateConsumerData(serviceName string, ndata utility.DataMap) {
	path, value := d.BindService(serviceName, ndata)
	zkClient.ZkCli.UpdateValue(path, value)
}

func init() {
	ServiceConsumer = &serviceConsumer{}
	ServiceConsumer.dataMap = utility.NewDataMap()
	ServiceConsumer.dataMap.Set("domain", zkClient.Domain)
	ServiceConsumer.PathList = make(map[string]string)
	ServiceConsumer.dataMap.Set("ip", zkClient.LocalIP)
}
