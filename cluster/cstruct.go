package cluster

import (
	"fmt"
	"github.com/colinyl/lib4go/utility"
	zk "github.com/colinyl/lib4go/zkClient"
	"github.com/colinyl/dataservicebus/config"
	"time"
	//"lib4go/logger"
	"log"
	"sync"
)

//-------------------------dsbcenter node value----------------------------

type dsbCenterNodeValue struct {
	IP       string
	IsMaster bool
	Last     int64
}

//ToString register center string
func (d *dsbCenterNodeValue) ToString() string {
	return fmt.Sprintf("IP:%s,IsMaster:%s,Last:%d",
		d.IP, d.IsMaster, d.Last)
}

//--------------------------------------------------------------------

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
	dataMap        utility.DataMap
	IsMasterServer bool
	Last           int64
}

//RegisterCenter  registercenter
var RegisterCenter *registerCenter

//--------------------------------------------------------------------

//-------------------------job consumer----------------------------
type jobConsumer struct {
	dataMap  utility.DataMap
	Last     int64
	PathList map[string]string
}

//JobConsumer service provider
var JobConsumer *jobConsumer

//--------------------------------------------------------------------

//-------------------------job manger----------------------------
type jobManager struct {
	dataMap  utility.DataMap
	Last     int64
	PathList map[string]string
	locker   sync.Mutex
}

//JobConfigItem job config item
type JobConfigItem struct {
	Name    string
	Script  string
	Trigger string
	Run     int
}

//JobConfigs job configs
type JobConfigs struct {
	Jobs map[string]*JobConfigItem
}

//JobManager job manager
var JobManager *jobManager

//CurrentJobConfigs current job configs
var CurrentJobConfigs *JobConfigs

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
}

//ServiceProvider service provider
var ServiceProvider *serviceProvider

//--------------------------------------------------------------------

//-------------------------service provider list----------------------------

var Log clusterLogger

type clusterLogger interface {
	Error(error)
	Info(s string)
	Infof(f string, arg ...interface{})
}

type defaultLogger struct {
}

func (d defaultLogger) Error(e error) {
	if e != nil {
		log.Fatal(e)
	}

}
func (d defaultLogger) Info(s string) {
	log.Print(s)
}
func (d defaultLogger) Infof(f string, arg ...interface{}) {
	log.Printf(f, arg...)
}

//--------------------------------------------------------------------

type zkClientObj struct {
	ZkCli   *zk.ZKCli
	LocalIP string
	Domain  string
	Err     error
}

var zkClient *zkClientObj

func init() {
	Log = &defaultLogger{}
	zkClient = &zkClientObj{}
	zkClient.Domain = config.Get().Domain
	zkClient.LocalIP = utility.GetLocalIP("192.168")
	zkClient.ZkCli, zkClient.Err = zk.New(config.Get().ZKServers, time.Second)
	Log.Error(zkClient.Err)
}
