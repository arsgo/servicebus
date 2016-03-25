package cluster

import (
	"lib4go/utility"
	"testing"
	"time"
)

func Test_dsbcenter(t *testing.T) {  
	 RegisterCenter.Bind()
	if !zkClient.ZkCli.Exists(RegisterCenter.Path) {
		t.Error("bind failed")
	}

	//test register center master
	dsb, s := RegisterCenter.GetMaskerValue()
	if dsb.IsMaster != RegisterCenter.IsMasterServer {
		t.Errorf("is not master server, please check zk server node:%s\r\n org:%s", dsb.ToString(), s)
	}

	spl, erx := RegisterCenter.DownloadServiceProviders()
	if erx != nil {
		t.Error(erx)
	}
	if len(spl) == 0 {
		ServiceProvider.BindServices([]string{"order.pay.get"}, utility.NewDataMap())
		spl, erx = RegisterCenter.DownloadServiceProviders()
		if erx != nil {
			t.Error(erx)
		}
	}
	if len(spl) == 0 {
		t.Error("bind service error")
	}

	for i, v := range spl {
		RegisterCenter.dataMap.Set("serviceName", i)
		for _, s := range v {
			RegisterCenter.dataMap.Set("ip", s)
			npath := RegisterCenter.dataMap.Translate(serviceProviderPath)
			err := zkClient.ZkCli.Delete(npath)
			if err != nil {
				t.Errorf("%s,path:%s", err, npath)
			}
		}
		npath := RegisterCenter.dataMap.Translate(serviceProviderRoot)
		err := zkClient.ZkCli.Delete(npath)
		if err != nil {
			t.Error(err)
		}
	}
	spl, erx = RegisterCenter.DownloadServiceProviders()
	if erx != nil {
		t.Error(erx)
	}
	if len(spl) != 0 {
		t.Error("delete service error")
	}

	//test service provider
	err := ServiceProvider.BindServices([]string{"order.pay.get"}, utility.NewDataMap())
	if err != nil {
		t.Error(err)
	}
	spl, erx = RegisterCenter.DownloadServiceProviders()
	if erx != nil {
		t.Error(erx)
	}
	if _, ok := spl["order.pay.get"]; !ok {
		t.Error("service bind error")
	}

	//test service consumer
	path, _ := ServiceConsumer.BindService("order.pay.get", ServiceConsumer.dataMap)
	if !zkClient.ZkCli.Exists(path) {
		t.Errorf("service consumer not create:%s", path)
	}

	//job
	err = JobManager.Bind()
	if err != nil {
		t.Error(err)
	}

	currentJobConfigs := &JobConfigs{}
	currentJobConfigs.Jobs = make(map[string]*JobConfigItem)
	currentJobConfigs.Jobs["auto.get.save.info.abc"] = &JobConfigItem{Name: "save user info", Script: "save.lua",
		Run: 0, Trigger: "10m"}
	err = JobManager.PublishJobConfigs(currentJobConfigs)
	if err != nil {
		t.Error(err)
	}
	currentConfigs, errm := JobManager.GetJobConfig()
	if errm != nil {
		t.Error(errm)
	}
	if len(currentConfigs.Jobs) != len(currentJobConfigs.Jobs) {
		t.Error("publish job configs error")
	}

	time.Sleep(time.Second * 3)
	if len(CurrentJobConfigs.Jobs) != len(currentJobConfigs.Jobs) {
		t.Error("job configs watcher error")
	}

	joblist, erm := JobManager.DownloadJobConsumers()
	if erm != nil {
		t.Error(erm)
	}
	count := len(joblist)
	err = JobConsumer.RegisterJobConsumer("auto.get.save.info.abc", utility.NewDataMap())
	if err != nil {
		t.Error(err)
		t.Errorf("count:%d", len(CurrentJobConfigs.Jobs))
	}
	joblist, erm = JobManager.DownloadJobConsumers()
	if erm != nil {
		t.Error(erm)
	}
	if len(joblist) != count+1 {
		t.Errorf("job consumer create error")
	}

	err = JobConsumer.RegisterJobConsumer("job.query.userinfo1232222322", utility.NewDataMap())
	if err != errjobNotConfig {
		t.Error("job consumer register error")
	}
	time.Sleep(time.Microsecond)
	//time.Sleep(time.Hour)
}
