package rpc

import (
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/colinyl/servicebus/rpc/rpcservice"
)

type ServiceProviderHandler interface {
	Request(name string, input string) (r string, err error)
	Send(name string, input string, data []byte) (r string, err error)
	GetSnap()(map[string]*serviceRequestSnap)
}

//JobProviderServer
type ServiceProviderServer struct {
	Address string
	Handler ServiceProviderHandler
}

//Serve
func (rpcServer *ServiceProviderServer) Serve() (er error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, err := thrift.NewTServerSocketTimeout(rpcServer.Address, time.Hour*24*31)
	if err != nil {
		return err
	}

	processor := rpcservice.NewServiceProviderProcessor(rpcServer.Handler)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	go server.Serve()
	er = nil
	return
}

//NewServiceProviderServer
func NewServiceProviderServer(address string, handler ServiceProviderHandler) *ServiceProviderServer {
	return &ServiceProviderServer{Address: address, Handler: handler}
}

type serviceRequestSnap struct {
	total   int
	current int
	success int
	failed  int
}

type ServiceHandler struct {
	services map[string]*serviceRequestSnap
}

func NewServiceHandler() *ServiceHandler {
	return &ServiceHandler{services: make(map[string]*serviceRequestSnap)}
}

func (s *ServiceHandler) Request(name string, input string) (r string, err error) {
	if _, ok := s.services[name]; !ok {
		s.services[name] = &serviceRequestSnap{}
	}
	s.services[name].total++
	s.services[name].current++
	defer func() {
		s.services[name].current--
		if err == nil {
			s.services[name].success++
		} else {
			s.services[name].failed++
		}
	}()
	return ServiceProviderPool.Request(name, input)
}

func (s *ServiceHandler) Send(name string, input string, data []byte) (r string, err error) {
	if _, ok := s.services[name]; !ok {
		s.services[name] = &serviceRequestSnap{}
	}
	s.services[name].total++
	s.services[name].current++
	defer func() {
		s.services[name].current--
		if err == nil {
			s.services[name].success++
		} else {
			s.services[name].failed++
		}
	}()
	return ServiceProviderPool.Send(name, input, data)
}

func (s *ServiceHandler) GetSnap() map[string]*serviceRequestSnap {
	return s.services
}
