package cluster

import (
	"encoding/json"
	"sort"

	"github.com/colinyl/lib4go/utility"
)

//Bind bind watcher for job configs changing
func (d *jobManager) Bind() error {
	config, err := d.GetConfigs()
	if err != nil {
		return err
	}
	CurrentJobConfigs = config
	return nil
}

//PublishJobConfig publish job configs
func (d *jobManager) PublishConfigs(configs *JobConfigs) error {
	data, err := json.Marshal(configs)
	if err != nil {
		return err
	}
	value := string(data)
	configPath := d.dataMap.Translate(jobConfigPath)
	err = zkClient.ZkCli.UpdateValue(configPath, value)
	return err
}

//GetConfigs Get job Config
func (d *jobManager) GetConfigs() (*JobConfigs, error) {
	defConfigs := &JobConfigs{}
	defConfigs.Jobs = make(map[string]*JobConfigItem)
	configPath := d.dataMap.Translate(jobConfigPath)
	value, err := zkClient.ZkCli.GetValue(configPath)
	if err != nil {
		Log.Info("get job config error")
		return defConfigs, err
	}
	err = json.Unmarshal([]byte(value), defConfigs)
	return defConfigs, err
}

//WatchConfigChange  watch job configs changes
func (d *jobManager) WatchConfigsChange(callback func(config *JobConfigs, err error)) {
	configPath := d.dataMap.Translate(jobConfigPath)
	jobChanges := make(chan string, 10)
	go zkClient.ZkCli.WatchValue(configPath, jobChanges)
	go func() {
		callback(d.GetConfigs())
		for {
			select {
			case <-jobChanges:
				{
					configs, err := d.GetConfigs()
					if err != nil {
						d.locker.Lock()
						callback(configs, err)
						d.locker.Unlock()
					}
				}
			}
		}
	}()
}

//DownloadJobConsumers   download job consumers
func (d *jobManager) DownloadConsumers() (JobConsumerList, error) {
	var jobList JobConsumerList = make(map[string][]string)
	serviceList, err := zkClient.ZkCli.GetChildren(d.dataMap.Translate(jobRoot))
	if err != nil {
		return jobList, err
	}

	for _, v := range serviceList {
		nmap := d.dataMap.Copy()
		nmap.Set("jobName", v)
		consumerList, er := zkClient.ZkCli.GetChildren(nmap.Translate(jobConsumerRoot))
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

func init() {
	JobManager = &jobManager{}
	CurrentJobConfigs := &JobConfigs{}
	CurrentJobConfigs.Jobs = make(map[string]*JobConfigItem)
	JobManager.dataMap = utility.NewDataMap()
	JobManager.dataMap.Set("domain", zkClient.Domain)
	JobManager.dataMap.Set("ip", zkClient.LocalIP)
	JobManager.PathList = make(map[string]string)
}
