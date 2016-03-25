package rpc

import (
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/colinyl/servicebus/rpc/rpcservice"
)

type ServiceProviderHandler interface {
	Request(name string, input string) (r string, err error)
    Send(name string, input string, data []byte) (r string, err error)
    
}

//JobProviderServer
type ServiceProviderServer struct {
	Address string
	Handler ServiceProviderHandler
}

//Serve
func (rpcServer *ServiceProviderServer) Serve() error {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, err := thrift.NewTServerSocketTimeout(rpcServer.Address, time.Hour*12)
	if err != nil {
		return err
	}

	processor := rpcservice.NewServiceProviderProcessor(rpcServer.Handler)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	server.Serve()
	return nil

}

//NewServiceProviderServer
func NewServiceProviderServer(address string, handler ServiceProviderHandler) *ServiceProviderServer {
	return &ServiceProviderServer{Address: address, Handler: handler}
}
