package rpc

import "log"

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	server := NewJobConsumerServer("127.0.0.1:1016", &serverHander{})
	go server.Serve()
}
