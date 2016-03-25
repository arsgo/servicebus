package rpc

import (
	"log"
)

type ServiceHandler struct {
}

func (s *ServiceHandler)Request(name string, input string) (r string, err error) {
   return "success",nil
}
func (s *ServiceHandler)Send(name string, input string, data []byte) (r string, err error) {
 return "success",nil
}

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	server := NewServiceProviderServer("127.0.0.1:1016", &ServiceHandler{})
	go server.Serve()

}
