package rpc

import (
	p "github.com/lib4go/lib4go/pool"
	"github.com/lib4go/dataservicebus/rpc/base/serviceprovider"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type serviceProviderClientObject struct {
	Address   string
	Port      string
	transport thrift.TTransport
	client    *serviceprovider.ServiceProviderClient
}

type serviceProviderClientPool struct {
	pool *p.ObjectPool
}

type serviceProviderClientFactory struct {
	ips   []string
	index int
}

func newserviceProviderClientObject(address string) *serviceProviderClientObject {
	return &serviceProviderClientObject{Address: address, Port: "1016"}
}
func (client *serviceProviderClientObject) open() error {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(client.Address)
	if err != nil {
		return err
	}

	useTransport := transportFactory.GetTransport(transport)
	newclient := serviceprovider.NewServiceProviderClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		return err
	}
	client.client = newclient
	client.transport = transport
	return nil
}

func (j *serviceProviderClientObject) request(jobName string, input string) (string, error) {
	return j.client.Request(jobName, input)
}

func (j *serviceProviderClientObject) Close() {
	j.transport.Close()
}

func newserviceProviderClientFactory(ips []string) *serviceProviderClientFactory {
	return &serviceProviderClientFactory{ips: ips}
}
func (j *serviceProviderClientFactory) Create() p.Object {
	v := j.index % len(j.ips)
	o := newserviceProviderClientObject(j.ips[v])
	o.open()
	return o
}
