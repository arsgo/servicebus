package cluster

import (
	//"log"
	"time"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
	zk "github.com/colinyl/lib4go/zkClient"
	"github.com/colinyl/servicebus/config"
	"github.com/colinyl/servicebus/rpc"
)

//-------------------------dsbcenter node value----------------------------

type dsbCenterNodeValue struct {
	IP     string
	Server string
	Last   string
}

//-------------------------service provider list----------------------------

//ServiceProviderList  service providers
type ServiceProviderList map[string][]string

//Add  add a service to list
func (s ServiceProviderList) Add(serviceName string, server string) {
	if s[serviceName] == nil {
		s[serviceName] = []string{}
	}
	s[serviceName] = append(s[serviceName], server)
}

//--------------------------------------------------------------------

//-------------------------register center----------------------------
type registerCenter struct {
	Path           string
	IP             string
	Port           string
	Server         string
	dataMap        utility.DataMap
	IsMasterServer bool
	Last           int64
	OnlineTime     int64
	LastPublish    int64
	jobCallback    func(config *JobConfigs, err error)
	Log            *logger.Logger
	rpcServer      *rpc.ServiceProviderServer
}

//RegisterCenter  registercenter
var RegisterCenter *registerCenter

//--------------------------------------------------------------------

//-------------------------job consumer----------------------------
type jobConsumer struct {
	dataMap  utility.DataMap
	Last     int64
	PathList map[string]string
	configs  *JobConfigs
	log      *logger.Logger
}

//JobConsumer service provider
var JobConsumer *jobConsumer

//--------------------------------------------------------------------

//JobConfigItem job config item
type JobConfigItem struct {
	Name        string
	Script      string
	Trigger     string
	Concurrency int
}

//JobConfigs job configs
type JobConfigs struct {
	Jobs map[string]JobConfigItem
}

//--------------------------------------------------------------------

//-------------------------job consumer list----------------------------

//JobConsumerList JobConsumerList
type JobConsumerList map[string][]string

//Add add a job consumer
func (s JobConsumerList) Add(jobName string, server string) {
	if s[jobName] == nil {
		s[jobName] = []string{}
	}
	s[jobName] = append(s[jobName], server)
}

//-------------------------service consumer list----------------------------

type serviceConsumer struct {
	dataMap  utility.DataMap
	Last     int64
	PathList map[string]string
}

//ServiceConsumer service consumer
var ServiceConsumer *serviceConsumer

//--------------------------------------------------------------------

//-------------------------service provider list----------------------------
type serviceProvider struct {
	Path    string
	dataMap utility.DataMap
	Last    int64
	Log     *logger.Logger
    Port    string
}

//ServiceProvider service provider
var ServiceProvider *serviceProvider

//--------------------------------------------------------------------

//-------------------------service provider list----------------------------

//--------------------------------------------------------------------

type zkClientObj struct {
	ZkCli   *zk.ZKCli
	LocalIP string
	Domain  string
	Err     error
}

var zkClient *zkClientObj

func init() {
	zkClient = &zkClientObj{}
	zkClient.Domain = config.Get().Domain
	zkClient.LocalIP = utility.GetLocalIP("192.168")
	zkClient.ZkCli, zkClient.Err = zk.New(config.Get().ZKServers, time.Second)
    
	//log.Print(zkClient.Err)
}
