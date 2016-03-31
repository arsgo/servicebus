package cluster

import (
	"encoding/json"
	"sort"
)

//PublishJobConfig publish job configs
func (d *registerCenter) PublishConfigs(configs *JobConfigs) error {
	data, err := json.Marshal(configs)
	if err != nil {
		return err
	}
	value := string(data)
	configPath := d.dataMap.Translate(jobConfigPath)
	err = zkClient.ZkCli.UpdateValue(configPath, value)
	return err
}

//DownloadJobConsumers   download job consumers
func (d *registerCenter) DownloadJobConsumers() (JobConsumerList, error) {
	var jobList JobConsumerList = make(map[string][]string)
	jobRootPath := d.dataMap.Translate(jobRoot)
	serviceList, err := zkClient.ZkCli.GetChildren(jobRootPath)	
	if err != nil {
		return jobList, err
	}

	for _, v := range serviceList {
		nmap := d.dataMap.Copy()
		nmap.Set("jobName", v)
        jobConsumerPath:=nmap.Translate(jobConsumerRoot)
		consumerList, er := zkClient.ZkCli.GetChildren(jobConsumerPath)
		if er != nil {
			return jobList, er
		}
		sort.Sort(sort.StringSlice(consumerList))
		for _, l := range consumerList {
			jobList.Add(v, l)
		}
	}
	return jobList, nil
}
