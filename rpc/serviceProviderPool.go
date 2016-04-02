package rpc

import (
	"fmt"
	"log"
	"time"

	"github.com/colinyl/lib4go/logger"
	p "github.com/colinyl/lib4go/pool"
	"github.com/colinyl/lib4go/utility"
)

//serviceProviderPool consumer connections  pool
var ServiceProviderPool *serviceProviderPool

//Register 注册服务列表
func (s *serviceProviderPool) Register(svs map[string][]string) {
	s.Log.Info("register service")
	nsvs := CrateServices(svs)
	//标记不能使用的服务
	for svName, names := range s.services.data {
		for ip := range names {
			if !nsvs.Contains(svName, ip) {
				nsvs.data[svName][ip].Status = false
			}
		}
	}

	//添加可以使用使用的服务
	for sv, servers := range nsvs.data {
		for ip := range servers {
			if !s.services.Contains(sv, ip) {
				s.Log.Info("add new service")
				s.pool.Register(ip, newserviceProviderClientFactory(ip), 1)
				s.services.Add(sv, ip)
			}
		}
	}
}

//UnRegister 取消注册服务列表
func (p *serviceProviderPool) UnRegister(svName string) {
	p.pool.UnRegister(svName)
}

//Request 执行Request请求
func (p *serviceProviderPool) Request(name string, input string) (string, error) {
	p.services.lk.Lock()
	defer p.services.lk.Unlock()

	if len(p.services.data[name]) == 0 {
		return p.getResult("500", "not find sp server"), nil
	}
	get := func(sv string) (result string, err error) {
        p.Log.Infof("request:%s",sv)
		o, err := p.pool.Get(sv)
		if err != nil {
			o.Fatal()
			return p.getResult("500", err.Error()), nil
		}
		if !o.Check() {
			return p.getResult("400", fmt.Sprintf("sp server not available:%s", sv)), nil
		}
		defer p.pool.Recycle(sv, o)
		obj := o.(*serviceProviderClient)
		return obj.Request(name, input)
	}

	for sv := range p.services.data[name] {
		result, err := get(sv)
		if err == nil {
			return result, nil
		}
	}
	return p.getResult("500", "not find available sp server"), nil

}

//Send 发送Send请求
func (p *serviceProviderPool) Send(name string, input string, data []byte) (result string, err error) {
	if len(p.services.data[name]) == 0 {
		return p.getResult("500", "not find sp server"), nil
	}
	get := func(sv string) (result string, err error) {
		o, err := p.pool.Get(sv)
		if err != nil {
			o.Fatal()
			return p.getResult("500", err.Error()), nil
		}
		if !o.Check() {
			return p.getResult("400", "sp server not available"), nil
		}
		defer p.pool.Recycle(sv, o)
		obj := o.(*serviceProviderClient)
		return obj.Send(name, input, data)
	}
	p.services.lk.Lock()
	defer p.services.lk.Unlock()
	for sv := range p.services.data[name] {
		result, err := get(sv)
		if err != nil {
			return result, nil
		}
	}
	return p.getResult("500", "not find available sp server"), nil

}
func (p *serviceProviderPool) clearUp() {
	timepk := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-timepk.C:
			{
				p.services.lk.Lock()
				for _, server := range p.services.data {
					for k, ip := range server {
						if !ip.Status && p.pool.Close(ip.IP) {
							delete(server, k)
						}
					}
				}
				p.services.lk.Unlock()
			}
		}
	}
}
func (p *serviceProviderPool) getResult(code string, msg string) string {
	dmap := utility.NewDataMap()
	dmap.Set("code", code)
	dmap.Set("msg", msg)
	return dmap.Translate(ERROR_FORMAT)
}

func NewServicePool() *serviceProviderPool {
	var err error
	pl := &serviceProviderPool{}
	pl.pool = p.New()
	pl.services = NewServices()
	pl.Log, err = logger.New("service privider pool", true)
	if err != nil {
		log.Println(err)
	}
	go pl.clearUp()
	return pl
}

func init() {
	ServiceProviderPool = NewServicePool()
}
