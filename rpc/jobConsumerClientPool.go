package rpc

import (
	p "github.com/colinyl/lib4go/pool"
)

//JobConsumerClientPool consumer connections  pool
var JobConsumerClientPool *jobConsumerClientPool

//Consumers job consumers map
type Consumers map[string][]string

func NewConsumers() Consumers {
	return make(map[string][]string)
}

func (c Consumers) Add(name string, ip string) {
	if _, ok := c[name]; !ok {
		c[name] = []string{}
	}
	c[name] = append(c[name], ip)
}
//Register register consumer list
func (p *jobConsumerClientPool) Register(consumers Consumers) {
	for i, v := range consumers {
		p.pool.Register(i, newJobConsumerClientFactory(v), 1)
	}
}

//Call request job consumer from pool
func (p *jobConsumerClientPool) Call(jobName string, input string) (string, error) {
	o, err := p.pool.Get(jobName)
	if err != nil {
		return "", err
	}
	defer p.pool.Recycle(jobName, o)
	obj := o.(*jobConsumerClientObject)
	return obj.call(jobName, input)
}

func init() {
	JobConsumerClientPool = &jobConsumerClientPool{}
	JobConsumerClientPool.pool = p.New()
}
