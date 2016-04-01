package main

import (
	//"time"

	"github.com/colinyl/servicebus/cluster"
)

func main() {

	cluster.ServiceProvider.StartRPC()

	cluster.ServiceProvider.WatchServiceConfigChange()
	//time.Sleep(time.Hour)

}
