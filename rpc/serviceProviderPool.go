package rpc

import (
	"time"

	p "github.com/colinyl/lib4go/pool"
)

//serviceProviderPool consumer connections  pool
var ServiceProviderPool *serviceProviderPool

//Register 注册服务列表
func (s *serviceProviderPool) Register(svs map[string][]string) {
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
func (p *serviceProviderPool) Request(name string, input string) (result string, err error) {
	get := func(sv string) (string, error) {
		o, err := p.pool.Get(sv)
		if err != nil {
			return "", err
		}
		if !o.Check() {
			return "", ERR_NOT_FIND_OBJ
		}
		defer p.pool.Recycle(sv, o)
		obj := o.(*serviceProviderClient)
		return obj.Request(name, input)
	}
	p.services.lk.Lock()
	defer p.services.lk.Unlock()
	for sv := range p.services.data[name] {
		result, err = get(sv)
		if err != nil {
			return
		}
	}
	return

}

//Send 发送Send请求
func (p *serviceProviderPool) Send(name string, input string, data []byte) (result string, err error) {
	get := func(sv string) (string, error) {
		o, err := p.pool.Get(sv)
		if err != nil {
			return "", err
		}
		defer p.pool.Recycle(sv, o)
		obj := o.(*serviceProviderClient)
		return obj.Send(name, input, data)
	}
	p.services.lk.Lock()
	defer p.services.lk.Unlock()
	for sv := range p.services.data[name] {
		result, err = get(sv)
		if err != nil {
			return
		}
	}
	return
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

func NewServicePool() *serviceProviderPool {
	pl := &serviceProviderPool{}
	pl.pool = p.New()
	pl.services = NewServices()
	go pl.clearUp()
	return pl
}

func init() {
	ServiceProviderPool = NewServicePool()
}
