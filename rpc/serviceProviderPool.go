package rpc

import (
	p "github.com/colinyl/lib4go/pool"
)

//serviceProviderPool consumer connections  pool
var ServiceProviderPool *serviceProviderPool

//Services 服务提供列表
type Services map[string][]string

func NewServices() Services {
	return make(map[string][]string)
}

func (c Services) Add(serviceName string, serverIP string) {
	if _, ok := c[serviceName]; !ok {
		c[serviceName] = []string{}
	}
	c[serviceName] = append(c[serviceName], serverIP)
}

//Register 注册服务列表
func (p *serviceProviderPool) Register(services Services) {
	p.addServices(services)
	for _, ips := range services {
		for _, v := range ips {
			p.pool.Register(v, newserviceProviderClientFactory(v), 1)
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

	for _, sv := range p.services[name] {
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

	for _, sv := range p.services[name] {
		result, err = get(sv)
		if err != nil {
			return
		}
	}
	return
}

func init() {
	ServiceProviderPool = &serviceProviderPool{}
	ServiceProviderPool.pool = p.New()
	ServiceProviderPool.services = make(map[string][]string)
}

func Register(services Services) {
    ServiceProviderPool.Register(services)
}

func UnRegister(svName string) {
    ServiceProviderPool.UnRegister(svName)
}

func  Request(name string, input string) (result string, err error) {
   return ServiceProviderPool.Request(name,input)
}

func  Send(name string, input string, data []byte) (result string, err error) {
    return ServiceProviderPool.Send(name,input,data)
}