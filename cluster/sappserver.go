package cluster

import (
	"encoding/json"
	"log"
	"time"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

type AutoConfig struct {
	Trigger string
	Script  string
}

type AppConfig struct {
	Status string
	Auto   []*AutoConfig
}

func (d *appServer) WatchConfigChange(callback func(config *AppConfig) error) {
	path := d.dataMap.Translate(appServerConfig)
START:
	if !zkClient.ZkCli.Exists(path) {
        d.Log.Infof("app config not exists:%s", path)
		time.Sleep(time.Second * 30)
		goto START
	}
	d.downloadConfig(callback)
	changes := make(chan string, 10)
	go zkClient.ZkCli.WatchValue(path, changes)
	for {
		select {
		case <-changes:
			{
				d.Log.Info("app config changed notices")
				d.downloadConfig(callback)
			}
		}
	}
}
func (d *appServer) downloadConfig(callback func(config *AppConfig) error) (config *AppConfig, err error) {
	config = &AppConfig{}
	path := d.dataMap.Translate(appServerConfig)
	values, err := zkClient.ZkCli.GetValue(path)
	if err != nil {
		d.Log.Error(err)
		return
	}
	err = json.Unmarshal([]byte(values), &config)
	if err != nil {
		d.Log.Error(err)
		return
	}
	err = callback(config)
	if err != nil {
		d.Log.Error(err)
		return
	}
	return
}

func init() {
	var err error
	AppServer = &appServer{}
	AppServer.dataMap = utility.NewDataMap()
	AppServer.dataMap.Set("domain", zkClient.Domain)
	AppServer.dataMap.Set("ip", zkClient.LocalIP)
	AppServer.Log, err = logger.New("app server", true)
	if err != nil {
		log.Print(err)
	}

}
