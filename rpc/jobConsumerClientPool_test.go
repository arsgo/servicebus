package rpc

import (
	"log"
	"strings"
	"testing"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	server := NewJobConsumerServer("127.0.0.1:1016", &serverHander{})
	go server.Serve()
}

type serverHander struct {
}

func (s *serverHander) Call(jobName string, input string) (r string, err error) {
	log.Printf("job:%s,input:%s", jobName, input)
	return "success", nil
}

func Test_rpc(t *testing.T) {
	jobc := NewConsumers()
	jobc.Add("job.user.get", "127.0.0.1:1016")
	JobConsumerClientPool.Register(jobc)
	o, err := JobConsumerClientPool.Call("job.user.get", "{}")
	if err != nil {
		t.Error(err)
	}
	if !strings.EqualFold(o, "success") {
		t.Error("return value error")
	}
}
