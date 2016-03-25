package main

import (
	//"fmt"
	"lib4go/utility"
    "time"
    "libs/cluster"
     "runtime"
    //"libs/zk"
)

var dataMap utility.DataMap
func init() {
	dataMap = utility.NewDataMap()
	dataMap.Set("domain", "/grs/delivery")
	dataMap.Set("ip", "127.0.0.1")
}



//批量下单，指定进程数，创建指定进程的程序，并批量进行下单请求
func main() {
    runtime.GOMAXPROCS(10)
	go cluster.RegisterCenter.Bind()
	
    

	/*time.Sleep(time.Second * 5)
	dsb, s := dsbcenter.RC.GetMasterValue()
	if dsb == nil || dsb.IsMaster != dsbcenter.RC.IsMasterServer || dsb.Last !=dsbcenter.RC.Last {
		fmt.Printf("is not master server, please check zk server node\r\n org:%s\r\n value:%s,master:[%s,%s],last:[%s,%s]\r\n", 
        s,dsb,dsb.IsMaster,dsbcenter.RC.IsMasterServer,dsb.Last,dsbcenter.RC.Last)
	}*/
	/*path := dsbcenter.RC.ClearServiceProviders()
	spl := dsbcenter.RC.DownloadServiceProviders()
	if len(spl) != 0 {
		fmt.Println("clear services  failed:%s", path)
	}
	dsbcenter.RC.BindService("order.pay.get")
	spl = dsbcenter.RC.DownloadServiceProviders()
	if len(spl) != 1 {
		fmt.Println("provider bind  failed")
	}*/
    time.Sleep(time.Hour)
}
