package rpc

import (
	p "github.com/colinyl/lib4go/pool"
)

//serviceProviderPool consumer connections  pool
var _serviceProviderPool *serviceProviderPool

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
func (p *serviceProviderPool) Register(Services services) {
	for i, v := range services {
		p.pool.Register(i, newserviceProviderClientFactory(v), 1)
	}
}

//Request 执行Request请求
func (p *serviceProviderPool) Request(name string, input string) (string, error) {
	o, err := p.pool.Get(name)
	if err != nil {
		return "", err
	}
	defer p.pool.Recycle(name, o)
	obj := o.(*serviceProviderObject)
	return obj.request(name, input)
}

//Send 发送Send请求
func (p *serviceProviderPool) Send(name string, input string, data []byte) (string, error) {
	o, err := p.pool.Get(name)
	if err != nil {
		return "", err
	}
	defer p.pool.Recycle(name, o)
	obj := o.(*serviceProviderObject)
	return obj.send(name, input, data)
}

func init() {
	_serviceProviderPool = &serviceProviderPool{}
	_serviceProviderPool.pool = p.New()
}
