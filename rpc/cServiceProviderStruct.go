package rpc

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	p "github.com/colinyl/lib4go/pool"
	"github.com/colinyl/servicebus/rpc/rpcservice"
)

type serviceProviderObject struct {
	Address   string
	Port      string
	transport thrift.TTransport
	client    *rpcservice.ServiceProviderClient
}

type serviceProviderPool struct {
	pool *p.ObjectPool
}

type serviceProviderClientFactory struct {
	ips   []string
	index int
}

func newserviceProviderObject(address string) *serviceProviderObject {
	return &serviceProviderObject{Address: address, Port: "1016"}
}

func (client *serviceProviderObject) open() error {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(client.Address)
	if err != nil {
		return err
	}

	useTransport := transportFactory.GetTransport(transport)
	newclient := rpcservice.NewServiceProviderClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		return err
	}
	client.client = newclient
	client.transport = transport
	return nil
}

func (j *serviceProviderObject) request(name string, input string) (string, error) {
	return j.client.Request(name, input)
}

func (j *serviceProviderObject) send(name string, input string, data []byte) (string, error) {
	return j.client.Send(name, input, data)
}

func (j *serviceProviderObject) Close() {
	j.transport.Close()
}

func newserviceProviderClientFactory(ips []string) *serviceProviderClientFactory {
	return &serviceProviderClientFactory{ips: ips}
}

func (j *serviceProviderClientFactory) Create() p.Object {
	v := j.index % len(j.ips)
	o := newserviceProviderObject(j.ips[v])
	o.open()
	return o
}
