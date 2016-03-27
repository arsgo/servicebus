package rpc

import (
	"log"
	"strings"
	"testing"
	"time"
)

type ServiceHandler struct {
}

func (s *ServiceHandler) Request(name string, input string) (r string, err error) {
	return name, nil
}

func (s *ServiceHandler) Send(name string, input string, data []byte) (r string, err error) {
	return input, nil
}

func Test_rpc(t *testing.T) {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	server := NewServiceProviderServer("192.168.1.106:8016", &ServiceHandler{})
	server.Serve()

	time.Sleep(time.Second)

	svs := NewServices()
	svs.Add("get_delivery_data", "192.168.1.106:8016")

	_serviceProviderPool.Register(svs)

	result, err := _serviceProviderPool.Request("get_delivery_data", "{}")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.EqualFold(result, "get_delivery_data") {
		t.Fatalf("result is error :%s", result)
	}

	result, err = _serviceProviderPool.Send("get_delivery_data", "{}", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.EqualFold(result, "{}") {
		t.Fatalf("result is error :%s", result)
	}

	//取消注册
	_serviceProviderPool.UnRegister("192.168.1.106:8016")
	result, err = _serviceProviderPool.Send("get_delivery_data", "{}", nil)
	if err == nil || !strings.EqualFold(err.Error(), "not find services: 192.168.1.106:8016") {
		t.Fatal("服务取消注册失败")
	}

}
