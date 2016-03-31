package scheduler

import (
	"encoding/json"
	"log"
	"time"
	//"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/servicebus/cluster"
	"github.com/colinyl/servicebus/rpc"
)

type JobDetail struct {
	config         *cluster.JobConfigItem
	funcDownClient func(name string) []string
}
type JobSnap struct {
	Name string
	Next time.Time
	Prev time.Time
}

func NewJob(config *cluster.JobConfigItem, funcDownClient func(name string) []string) *JobDetail {
	return &JobDetail{config: config, funcDownClient: funcDownClient}
}

func (j *JobDetail) Run() {
	log.Println("job trigger")
	var count int
	content, _ := json.Marshal(j.config)
	servers := j.funcDownClient(j.config.Name)
	for i := 0; i < len(servers); i++ {
		client := rpc.NewServiceProviderClient(servers[i])
		if client.Open() != nil {
			continue
		}
		_, err := client.Request(j.config.Name, string(content))
		if err != nil {
			count++
		}
		if count >= j.config.Concurrency {
			break
		}
	}
	if count != j.config.Concurrency {
	//	log.Printf("job run exception:%s\r\n", j.config.Name)
	}
}
