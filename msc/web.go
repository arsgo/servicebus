package main

import (
	"github.com/hoisie/web"
)

type monitor interface {
	GetSnap() string
}
type SysMonitor struct {
	monitors map[string]monitor
}

func NewSysMonitor() *SysMonitor {
	return &SysMonitor{}
}
func (m *SysMonitor) Add(name string, mo monitor) {
	m.monitors[name] = mo
}

func (m *SysMonitor) Start() {
	web.Get("/(.*)", m.show)
	web.Run("0.0.0.0:9999")
}

func (m *SysMonitor) show(method string)(string) {
     if v,ok:=m.monitors[method];ok{
        return v.GetSnap()
     }else{
         return method
     }
}
