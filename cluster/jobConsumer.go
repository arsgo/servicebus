package cluster

import (
	"fmt"
	"lib4go/utility"
	"time"
)

//UpdateServiceData  update consumer data
func (d *jobConsumer) UpdateJobData(jobName string, ndata utility.DataMap) error {
	nmap := d.dataMap.Merge(ndata)
	nmap.Set("jobName", jobName)
	nmap.Set("now", fmt.Sprintf("%d", time.Now().Unix()))
	path := nmap.Translate(jobConsumerPath)
    var (p string
        ok bool)
	if p, ok = d.PathList[path]; !ok {
		return errJobConsumerNotRegister
	} 
	return zkClient.ZkCli.UpdateValue(p, nmap.Translate(jobConsumerValue))	
}


//RegisterJobConsumer register job consumer
func (d *jobConsumer) RegisterJobConsumer(jobName string, ndata utility.DataMap) error {	
    if _,ok:=CurrentJobConfigs.Jobs[jobName];!ok{
        return errjobNotConfig
    }
    nmap := d.dataMap.Merge(ndata)
	nmap.Set("jobName", jobName)
	nmap.Set("now", fmt.Sprintf("%d", time.Now().Unix()))
	path := nmap.Translate(jobConsumerPath)
	if _, ok := d.PathList[path]; !ok {
		rpath, err := zkClient.ZkCli.CreateSeqNode(path, nmap.Translate(jobConsumerValue))
		if err != nil {
			return err
		}
		d.PathList[path] = rpath
	}
	return nil
}


func init() {
	JobConsumer = &jobConsumer{}
	JobConsumer.dataMap = utility.NewDataMap()
	JobConsumer.dataMap.Set("domain", zkClient.Domain)
    JobConsumer.dataMap.Set("ip", zkClient.LocalIP)
	JobConsumer.PathList = make(map[string]string)	
}