package cluster

import (
	"github.com/colinyl/servicebus/rpc"
)

func (d *registerCenter) StartRPC() {
	address := rpc.GetLocalRandomAddress()
	d.Port = address
    d.dataMap.Set("port",d.Port)
    d.updateNodeValue()
	rpcServer := rpc.NewServiceProviderServer(address, d.log, &rpc.ServiceHandler{})
	rpcServer.Serve()
}
