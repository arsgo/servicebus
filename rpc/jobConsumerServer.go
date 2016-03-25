package rpc

import (
	"fmt"
	"github.com/colinyl/dataservicebus/rpc/base/jobconsumer"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type JobConsumerHandler interface {
	Call(jobName string, input string) (r string, err error)
}

//JobConsumerServer 
type JobConsumerServer struct {
	Address string
	Handler JobConsumerHandler
}

//Serve
func (rpcServer *JobConsumerServer) Serve() error {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, err := thrift.NewTServerSocketTimeout(rpcServer.Address, time.Hour*12)
	if err != nil {
		return err
	}

	processor := jobconsumer.NewJobConsumerProcessor(rpcServer.Handler)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	fmt.Println("server in", rpcServer.Address)
	server.Serve()
	return nil
}

//NewJobConsumerServer
func NewJobConsumerServer(address string, handler JobConsumerHandler) *JobConsumerServer {
	return &JobConsumerServer{Address: address, Handler: handler}
}
