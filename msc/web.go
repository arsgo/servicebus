package main

import (
	"github.com/colinyl/web"
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
func (m *SysMonitor) Add(name string, m monitor) {
	m.monitors[name] = m
}

func (m *SysMonitor) Start() {
	web.Get("/(.*)", m.show)
	web.Run("0.0.0.0:9999")
}

func (m *SysMonitor) show(ctx *web.Context, val string) {
     key:=ctx.Params["method"]
     if v,ok:=m.monitors[key];ok{
         println[v.GetSnap()]
     }
}
