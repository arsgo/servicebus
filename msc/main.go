package main

import (
	"fmt"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
	"github.com/colinyl/servicebus/cluster"
	"github.com/colinyl/servicebus/rpc"
	"github.com/colinyl/servicebus/scheduler"
)

var log, _ = logger.New("message service center")

func main() {

	mnt := NewSysMonitor()
	//绑定注册中心,选举master,自动监控服务列表变化并重新发布服务信息,监控最新的服务列表,并重新注册POOL
	cluster.RegisterCenter.Bind()
	cluster.RegisterCenter.WatchServiceListChange(func(services map[string][]string, err error) {
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		var cservice rpc.Services = services
		rpc.ServiceProviderPool.Register(cservice)
	})
	mnt.Add("center", cluster.RegisterCenter)

	// 启动本地服务,等待service consumer远程连接,并转发到远程service provider
	rpcServer := rpc.NewServiceProviderServer(utility.GetLocalIP("192"), &rpc.ServiceHandler{})
	rpcServer.Serve()
	mnt.Add("server", rpcServer)

	//绑定JOB服务,并监控服务列表变化,并启动scheduler,根据cron配置时间远程调用job consumer
	cluster.JobManager.Bind()
	cluster.JobManager.WatchConfigsChange(func(config *cluster.JobConfigs, err error) {
		if len(config.Jobs) == 0 {
			return
		}
		scheduler.Stop()
		for _, v := range config.Jobs {
			scheduler.AddJob(scheduler.NewJob(v, func(name string) []string {
				consumers, er := cluster.JobManager.DownloadConsumers()
				if er != nil {
					log.Fatalf("下载JOB consumers error:\r\n", er.Error())
					return []string{}
				}
				return consumers[name]
			}))
		}
		scheduler.Start()
	})
	mnt.Add("job", cluster.JobManager)
	mnt.Add("scheduler", scheduler.Schd)
	 mnt.Start()
	fmt.Println("启动成功......")

}
