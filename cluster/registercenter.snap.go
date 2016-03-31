package cluster

import (
	"time"
    "fmt"
)

func (d *registerCenter) StartSnap() {
	timeTicker := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-timeTicker.C:
			d.Last = time.Now().Unix()
            d.dataMap.Set("last",fmt.Sprintf("%d",d.Last))
			d.Log.Infof("update snapshot:%d", d.Last)
			d.updateNodeValue()
		}
	}

}
