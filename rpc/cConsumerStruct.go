package rpc

import (
	p "github.com/colinyl/lib4go/pool"
	"github.com/colinyl/dataservicebus/rpc/base/jobconsumer"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type jobConsumerClientObject struct {
	Address   string
	Port      string
	transport thrift.TTransport
	client    *jobconsumer.JobConsumerClient
}

type jobConsumerClientPool struct {
	pool *p.ObjectPool
}

type jobConsumerClientFactory struct {
	ips   []string
	index int
}

func newJobConsumerClientObject(address string) *jobConsumerClientObject {
	return &jobConsumerClientObject{Address: address, Port: "1016"}
}
func (client *jobConsumerClientObject) open() error {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(client.Address)
	if err != nil {
		return err
	}

	useTransport := transportFactory.GetTransport(transport)
	newclient := jobconsumer.NewJobConsumerClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		return err
	}
	client.client = newclient
	client.transport = transport
	return nil
}

func (j *jobConsumerClientObject) call(jobName string, input string) (string, error) {
	return j.client.Call(jobName, input)
}

func (j *jobConsumerClientObject) Close() {
	j.transport.Close()
}

func newJobConsumerClientFactory(ips []string) *jobConsumerClientFactory {
	return &jobConsumerClientFactory{ips: ips}
}
func (j *jobConsumerClientFactory) Create() p.Object {
	v := j.index % len(j.ips)
	o := newJobConsumerClientObject(j.ips[v])
	o.open()
	return o
}
