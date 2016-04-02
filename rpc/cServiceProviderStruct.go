package rpc

import (
	"log"
	"net"
	"strings"
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/colinyl/lib4go/logger"
	p "github.com/colinyl/lib4go/pool"
	"github.com/colinyl/servicebus/rpc/rpcservice"
)

var (
    ERROR_FORMAT    =`{"code":"@code","msg":"@msg"}`
)

type serviceProviderClient struct {
	Address   string
	transport thrift.TTransport
	client    *rpcservice.ServiceProviderClient
	isFatal   bool
}
type rcServerService struct {
	Status bool
	IP     string
}

type rcServerPool struct {
	pool    *p.ObjectPool
	servers map[string]*rcServerService
	lk      sync.Mutex
	Log     *logger.Logger
}

type serviceProviderPool struct {
	pool     *p.ObjectPool
	services *Services
    Log     *logger.Logger
}

type serviceProviderClientFactory struct {
	ip string
}

type serviceItem struct {
	IP     string
	Status bool
}

type serviceCollection map[string]*serviceItem

//Services 服务提供列表
type Services struct {
	data map[string]serviceCollection
	lk   sync.Mutex
}

func NewServices() *Services {
	return &Services{data: make(map[string]serviceCollection)}
}

func CrateServices(sv map[string][]string) (s *Services) {
	s = NewServices()
	for k, c := range sv {
		if _, ok := s.data[k]; !ok {
			s.data[k] = make(map[string]*serviceItem)
		}

		for _, ip := range c {
			if _, ok := s.data[k][ip]; !ok {
				s.data[k][ip] = &serviceItem{IP: ip, Status: true}
			}
		}
	}
	return
}

func (c *Services) Contains(serviceName string, serverIP string) bool {
	c.lk.Lock()
	defer c.lk.Unlock()
	if _, ok := c.data[serviceName]; !ok {
		return false
	}
	if _, ok := c.data[serviceName][serverIP]; !ok {
		return false
	}
	return true
}

func (c *Services) Add(serviceName string, serverIP string) {
	c.lk.Lock()
	defer c.lk.Unlock()
	if _, ok := c.data[serviceName]; !ok {
		c.data[serviceName] = make(map[string]*serviceItem)
	}

	if _, ok := c.data[serviceName][serverIP]; !ok {
		c.data[serviceName][serverIP] = &serviceItem{IP: serverIP, Status: true}
	}
}

func (c *Services) Remove(serviceName string, serverIP string) {
	c.lk.Lock()
	defer c.lk.Unlock()
	if _, ok := c.data[serviceName]; !ok {
		return
	}

	if _, ok := c.data[serviceName][serverIP]; !ok {
		return
	}
	c.data[serviceName][serverIP].Status = false
}

func NewServiceProviderClient(address string) *serviceProviderClient {
	addr := address
	if !strings.Contains(address, ":") {
		addr = net.JoinHostPort(address, "1016")
	}
	return &serviceProviderClient{Address: addr}
}

func (client *serviceProviderClient) Open() (err error) {
	client.transport, err = thrift.NewTSocket(client.Address)
	if err != nil {
		return err
	}

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	pf := thrift.NewTBinaryProtocolFactoryDefault()

	useTransport := transportFactory.GetTransport(client.transport)
	client.client = rpcservice.NewServiceProviderClientFactory(useTransport, pf)
	if err := client.client.Transport.Open(); err != nil {
		return err
	}
	return nil
}

func (j *serviceProviderClient) Request(name string, input string) (string, error) {
	return j.client.Request(name, input)
}

func (j *serviceProviderClient) Send(name string, input string, data []byte) (string, error) {
	return j.client.Send(name, input, data)
}

func (j *serviceProviderClient) Close() {
	j.transport.Close()
}

func (j *serviceProviderClient) Check() bool {
	return !j.isFatal && j.transport != nil
}
func (j *serviceProviderClient) Fatal() {
	j.isFatal = true
}

func newserviceProviderClientFactory(ip string) *serviceProviderClientFactory {
	log.Println(ip)
	return &serviceProviderClientFactory{ip: ip}
}

func (j *serviceProviderClientFactory) Create() (p.Object, error) {
	o := NewServiceProviderClient(j.ip)
	if err := o.Open(); err != nil {
		return nil, err
	}
	return o, nil
}
