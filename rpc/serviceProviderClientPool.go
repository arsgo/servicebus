package rpc

import (
	p "github.com/lib4go/lib4go/pool"
)

//ServiceProviderClientPool consumer connections  pool
var ServiceProviderClientPool *serviceProviderClientPool

//Servers job Servers map
type Servers map[string][]string

func NewServers() Servers {
	return make(map[string][]string)
}

func (c Servers) Add(name string, ip string) {
	if _, ok := c[name]; !ok {
		c[name] = []string{}
	}
	c[name] = append(c[name], ip)
}

//Register register consumer list
func (p *serviceProviderClientPool) Register(Servers Servers) {
	for i, v := range Servers {
		p.pool.Register(i, newserviceProviderClientFactory(v), 1)
	}
}

//Call request job consumer from pool
func (p *serviceProviderClientPool) Call(jobName string, input string) (string, error) {
	o, err := p.pool.Get(jobName)
	if err != nil {
		return "", err
	}
	defer p.pool.Recycle(jobName, o)
	obj := o.(*serviceProviderClientObject)
	return obj.request(jobName, input)
}
func init() {
	ServiceProviderClientPool = &serviceProviderClientPool{}
	ServiceProviderClientPool.pool = p.New()
}
