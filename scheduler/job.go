package scheduler

import (
	"encoding/json"
	"strings"
"time"
	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/servicebus/cluster"
	"github.com/colinyl/servicebus/rpc"
)

var (
	log *logger.Logger
)

func init() {
	log, _ = logger.New("scheduler")
}

type JobDetail struct {
	config         *cluster.JobConfigItem
	funcDownClient func(name string) []string
}
type JobSnap struct{
    Name string
    Next time.Time
    Prev time.Time
}


func NewJob(config *cluster.JobConfigItem, funcDownClient func(name string) []string) *JobDetail {
	return &JobDetail{config: config, funcDownClient: funcDownClient}
}

func (j *JobDetail) Run() {

	log.Infof("run job:%s\r\n", j.config.Name)
	var index int
	content, _ := json.Marshal(j.config)
	servers := j.funcDownClient(j.config.Name)
	for i := 0; i < len(servers); i++ {
		client := rpc.NewServiceProviderClient(servers[i])
		if client.Open() != nil {
			continue
		}
		_, err := client.Request(j.config.Name, string(content))
		if err != nil {
			log.Infof("\trun job:%s:%s", j.config.Name, servers[i])
			index++
		}
		if index >= j.config.Count {
			break
		}
	}
	if index != j.config.Count {
		log.Errorf("job run exception:%s,%s\r\n", j.config.Name, strings.Join(servers, ","))
	}
}
