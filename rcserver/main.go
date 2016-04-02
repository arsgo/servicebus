package main

import (
	"runtime"

	"github.com/colinyl/servicebus/cluster"
	"github.com/colinyl/servicebus/rpc"
	"github.com/colinyl/servicebus/scheduler"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	cluster.RegisterCenter.Bind()
	cluster.RegisterCenter.WatchServicesChange(func(services map[string][]string, err error) {
		if err != nil {
			cluster.RegisterCenter.Log.Fatal(err.Error())
			return
		}
		rpc.ServiceProviderPool.Register(services)
	})

	cluster.RegisterCenter.WatchJobChange(func(config *cluster.JobConfigs, err error) {
		scheduler.Stop()
		if len(config.Jobs) == 0 {
			return
		}
		for _, v := range config.Jobs {
			scheduler.AddJob(scheduler.NewJob(&v, func(name string) []string {
				consumers, er := cluster.RegisterCenter.DownloadJobConsumers()
				if er != nil {
					cluster.RegisterCenter.Log.Fatalf("download JOB consumers error:%+v\r\n", er)
					return []string{}
				}
				return consumers[name]
			}))
		}
		scheduler.Start()
	})

	cluster.RegisterCenter.StartRPC()

	cluster.RegisterCenter.StartSnap()
}
