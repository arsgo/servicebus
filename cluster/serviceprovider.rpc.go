package cluster

import (
	"errors"
	"log"
	"strings"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/lua"
	"github.com/colinyl/servicebus/rpc"
)

type scriptEngine struct {
	script   *lua.LuaPool
	provider *serviceProvider
	Log      *logger.Logger
}

func (s *scriptEngine) Request(cmd string, input string) (string, error) {
    s.Log.Infof("request:%s",cmd)
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
	var err error
	en := &scriptEngine{script: lua.NewLuaPool(), provider: p}
	en.Log, err = logger.New("app script", true)
	if err != nil {
		log.Println(err)
	}
	return en
}

func (d *serviceProvider) StartRPC() {
	address := rpc.GetLocalRandomAddress()
	d.Port = address
	d.dataMap.Set("port", d.Port)
	rpcServer := rpc.NewServiceProviderServer(address, d.Log, rpc.NewServiceHandler(NewScript(d)))
	rpcServer.Serve()
}
