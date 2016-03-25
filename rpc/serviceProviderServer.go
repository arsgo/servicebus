package rpc

import (
	"github.com/lib4go/dataservicebus/rpc/base/serviceprovider"
	"time"
	"git.apache.org/thrift.git/lib/go/thrift"
)

type JobProviderHandler interface {
	Request(jobName string, input string) (r string, err error)
}

//JobProviderServer
type JobProviderServer struct {
	Address string
	Handler JobProviderHandler
}

//Serve
func (rpcServer *JobProviderServer) Serve() error {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, err := thrift.NewTServerSocketTimeout(rpcServer.Address, time.Hour*12)
	if err != nil {
		return err
	}

	processor := serviceprovider.NewServiceProviderProcessor(rpcServer.Handler)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	server.Serve()
	return nil

}

//NewServiceProviderServer
func NewServiceProviderServer(address string, handler JobProviderHandler) *JobProviderServer {
	return &JobProviderServer{Address: address, Handler: handler}
}
