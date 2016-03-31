package cluster

import (
	"strings"
	"testing"
	//"time"
	//"github.com/colinyl/lib4go/utility"
)
func Test_dsbcenter(t *testing.T) {

	master := RegisterCenter
	master.Bind()
	if !zkClient.ZkCli.Exists(master.Path) {
		t.Error("bind failed")
	}
  /*
	//test register center master
	masterValue, s := master.GetValue()
	if !strings.EqualFold(masterValue.Server, "master") || !master.IsMasterServer {
		t.Errorf("is not master server, please check zk server node:%s,%s,%v", s,
			masterValue.Server, master.IsMasterServer)
	}*/

	slave := NewRegisterCenter("/grs/delivery", "192.168.0.245")
	slave.Bind()
	slaveValue, st:= slave.getValue()
    if !strings.EqualFold(slaveValue.Server, "salve") || slave.IsMasterServer {
		t.Errorf("is not slave server, please check zk server node:%s,%s,%v", st,
			slaveValue.Server, slave.IsMasterServer)
	}
  /*
	spl, erx := master.DownloadServiceProviders()
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
		Count: 0, Trigger: "10m"}
	err = JobManager.PublishConfigs(currentJobConfigs)
	if err != nil {
		t.Error(err)
	}
	currentConfigs, errm := JobManager.GetConfigs()
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

	joblist, erm := JobManager.DownloadConsumers()
	if erm != nil {
		t.Error(erm)
	}
	count := len(joblist)
	err = JobConsumer.Register("auto.get.save.info.abc", utility.NewDataMap())
	if err != nil {
		t.Error(err)
		t.Errorf("count:%d", len(CurrentJobConfigs.Jobs))
	}
	joblist, erm = JobManager.DownloadConsumers()
	if erm != nil {
		t.Error(erm)
	}
	if len(joblist) != count+1 {
		t.Errorf("job consumer create error")
	}

	err = JobConsumer.Register("job.query.userinfo1232222322", utility.NewDataMap())
	if err != errjobNotConfig {
		t.Error("job consumer register error")
	}
     */
	//	time.Sleep(time.Second*10)
	//	time.Sleep(time.Hour)
       
}
