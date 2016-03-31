package cluster

import (
	"github.com/colinyl/servicebus/rpc"
)

func (d *registerCenter) StartRPC() {
	address := rpc.GetLocalRandomAddress()
	d.Port = address
	d.dataMap.Set("port", d.Port)
	d.updateNodeValue()
	d.rpcServer = rpc.NewServiceProviderServer(address, d.Log, &rpc.ServiceHandler{})
	d.rpcServer.Serve()
}
