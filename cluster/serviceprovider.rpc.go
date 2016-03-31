package cluster

import (
	"github.com/colinyl/servicebus/rpc"
)

func (d *serviceProvider) StartRPC() {
	address := rpc.GetLocalRandomAddress()
	d.Port = address
	d.dataMap.Set("port", d.Port)
	//d.updateNodeValue()
	rpcServer := rpc.NewServiceProviderServer(address, d.Log, &rpc.ServiceHandler{})
	rpcServer.Serve()
}
