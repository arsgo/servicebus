package cluster

import (
	"errors"
	"strings"

	"github.com/colinyl/lib4go/lua"
	"github.com/colinyl/servicebus/rpc"
)

type scriptEngine struct {
	script   *lua.LuaPool
	provider *serviceProvider
}

func (s *scriptEngine) Request(cmd string, input string) (string, error) {
	s.provider.lk.Lock()
	svs, ok := s.provider.services.services[cmd]
	s.provider.lk.Unlock()
	if !ok {
		return "", errors.New("not suport")
	}
	path := svs.Script
	values, err := s.script.Call(path, input)
	return strings.Join(values, ","), err
}
func (s *scriptEngine) Send(cmd string, input string, data []byte) (string, error) {
	s.provider.lk.Lock()
	svs, ok := s.provider.services.services[cmd]
	s.provider.lk.Unlock()

	if !ok {
		return "", errors.New("not suport")
	}
	path := svs.Script
	values, err := s.script.Call(path, input)
	return strings.Join(values, ","), err
}

func NewScript(p *serviceProvider) *scriptEngine {
	return &scriptEngine{script: lua.NewLuaPool(), provider: p}
}

func (d *serviceProvider) StartRPC() {
	address := rpc.GetLocalRandomAddress()
	d.Port = address
	d.dataMap.Set("port", d.Port)	
	rpcServer := rpc.NewServiceProviderServer(address, d.Log, rpc.NewServiceHandler(NewScript(d)))
	rpcServer.Serve()
}
