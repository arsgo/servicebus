package rpc

import (
	//"encoding/json"
	"fmt"
	"github.com/colinyl/lib4go/logger"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/colinyl/lib4go/net"
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
    log *logger.Logger
}

//Serve
func (rpcServer *ServiceProviderServer) Serve() (er error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, er := thrift.NewTServerSocketTimeout(rpcServer.Address, time.Hour*24*31)
	if er != nil {
		rpcServer.log.Error(er)
		return
	}

	processor := rpcservice.NewServiceProviderProcessor(rpcServer.Handler)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	rpcServer.log.Infof("::start rpc server :%s", rpcServer.Address)
	go func() {
		er = server.Serve()
		if er != nil {
			rpcServer.log.Error(er)
		}
	}()
	return
}

//NewServiceProviderServer
func NewServiceProviderServer(address string,log *logger.Logger, handler ServiceProviderHandler) *ServiceProviderServer {
	return &ServiceProviderServer{Address: address, Handler: handler,log:log}
}

func GetLocalRandomAddress() string {
	return fmt.Sprintf(":%d", getPort())
}

func getPort() int {
	for i := 0; i < 100; i++ {
		port := 1016 + i*8
		if net.IsTCPPortAvailable(port) {
			return port
		}
	}
	return -1
}
