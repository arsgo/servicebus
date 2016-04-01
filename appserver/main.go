package main

import (
	"github.com/colinyl/servicebus/cluster"
	"github.com/colinyl/servicebus/scheduler"
)

func main() {

	cluster.AppServer.WatchConfigChange(func(config *cluster.AppConfig) error {
		cluster.AppServer.Log.Info("::start server")
		scheduler.Stop()
		if len(config.Auto) == 0 {
			return nil
		}
		for _, v := range config.Auto {
			cluster.AppServer.Log.Info(v.Trigger)
			scheduler.AddTask(v.Trigger, scheduler.NewTask(v.Script, func(name string) {
				cluster.AppServer.Log.Infof("start_script:%s", name)
			}))
		}
		scheduler.Start()
		return nil
	})
}
