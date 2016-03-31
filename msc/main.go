package main

import (	
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/servicebus/cluster"
	"github.com/colinyl/servicebus/scheduler"
)

var vlog, _ = logger.New("message service center", true)
func main() {       
	cluster.RegisterCenter.Bind(vlog)
	cluster.RegisterCenter.WatchServicesChange(func(services map[string][]string, err error) {
		if err != nil {
			vlog.Fatal(err.Error())
			return
		}
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
					vlog.Fatalf("download JOB consumers error:%+v\r\n", er)
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
