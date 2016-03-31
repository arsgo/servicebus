package cluster

import (
	"encoding/json"
	"fmt"
	"time"

	//"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

//UpdateServiceData  update consumer data
func (d *jobConsumer) UpdateJobData(jobName string, ndata utility.DataMap) error {
	nmap := d.dataMap.Merge(ndata)
	nmap.Set("jobName", jobName)
	nmap.Set("now", fmt.Sprintf("%d", time.Now().Unix()))
	path := nmap.Translate(jobConsumerPath)
	var (
		p  string
		ok bool
	)
	if p, ok = d.PathList[path]; !ok {
		return errJobConsumerNotRegister
	}
	return zkClient.ZkCli.UpdateValue(p, nmap.Translate(jobConsumerValue))
}

//Register register job consumer
func (d *jobConsumer) Register(jobName string, ndata utility.DataMap) error {
	if _, ok := d.configs.Jobs[jobName]; !ok {
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

func (d *jobConsumer) getConfigs() (*JobConfigs, error) {
	defConfigs := &JobConfigs{}
	defConfigs.Jobs = make(map[string]JobConfigItem)
	configPath := d.dataMap.Translate(jobConfigPath)
	value, err := zkClient.ZkCli.GetValue(configPath)
	if err != nil {
		d.log.Info("get job config error")
		return defConfigs, err
	}
	err = json.Unmarshal([]byte(value), defConfigs)
	return defConfigs, err
}

//NewJobConsumer create a new job consumer
func NewJobConsumer() *jobConsumer {
	var err error
	jc := &jobConsumer{}
	jc.dataMap = utility.NewDataMap()
	jc.dataMap.Set("domain", zkClient.Domain)
	jc.dataMap.Set("ip", zkClient.LocalIP)
	jc.PathList = make(map[string]string)
	if err != nil {
		fmt.Println(err)
	}
	return jc
}

func init() {
	JobConsumer = NewJobConsumer()
}
