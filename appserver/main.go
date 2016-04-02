package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/colinyl/servicebus/cluster"
	"github.com/colinyl/servicebus/rpc"
	"github.com/colinyl/servicebus/scheduler"
)

func main() {

	log.SetFlags(log.Ldate | log.Lmicroseconds)
	///log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	//监控RC服务器变化
	cluster.AppServer.WatchRCServerChange(func(configs []*cluster.RCServerConfig) {
		cluster.AppServer.Log.Info("recv rc server changes")
		servers := make(map[string]string)
		for _, v := range configs {
			sv := fmt.Sprintf("%s%s", v.IP, v.Port)
			servers[sv] = sv
		}
		cluster.AppServer.Log.Infof("app servers:%d", len(servers))
		rpc.RCServerPool.Register(servers)
	})

	// 监控当前机器配置的变化
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
				rtvalues, err := cluster.AppLuaEngine.Call(name, v.Input)
				if err != nil {
					cluster.AppServer.Log.Error(err)
				} else {
					cluster.AppServer.Log.Infof("app result:%s", strings.Join(rtvalues, ","))
				}
			}))
		}
		scheduler.Start()
		return nil
	})

}
