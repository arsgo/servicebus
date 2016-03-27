package rpc

import (
	"errors"
	"net"
	"strings"
    "log"
	"git.apache.org/thrift.git/lib/go/thrift"
	p "github.com/colinyl/lib4go/pool"
	"github.com/colinyl/servicebus/rpc/rpcservice"
)

var (
	ERR_NOT_FIND_OBJ        error = errors.New("未找到可用的服务器连接")
	ERR_CANT_CONNECT_SERVER error = errors.New("未找到可用的服务器连接")
)

type serviceProviderClient struct {
	Address   string
	transport thrift.TTransport
	client    *rpcservice.ServiceProviderClient
	isFatal   bool
}

type serviceProviderPool struct {
	pool     *p.ObjectPool
	services map[string][]string
}

type serviceProviderClientFactory struct {
	ip    string
}

func (s *serviceProviderPool) addServices(svs Services) {
	for name, server := range svs {
		for _, sv := range server {
			s.services[name] = append(s.services[name], sv)
		}
	}
}

func NewServiceProviderClient(address string) *serviceProviderClient {
	addr:=address
	if !strings.Contains(address, ":") {
		addr=net.JoinHostPort(address,"1016")
	}
	return &serviceProviderClient{Address:addr}
}

func (client *serviceProviderClient) Open()(err error) {
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
func (j *serviceProviderClient) RequestFatal() {
	j.isFatal = true
}

func newserviceProviderClientFactory(ip string) *serviceProviderClientFactory {
    log.Println(ip)
	return &serviceProviderClientFactory{ip: ip}
}

func (j *serviceProviderClientFactory) Create() (p.Object, error) {
	o := NewServiceProviderClient(j.ip)
	if err:=o.Open();err != nil {
		return nil, err
	}
	return o,nil
}
