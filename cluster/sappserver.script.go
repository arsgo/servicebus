package cluster

import (
	"github.com/colinyl/lib4go/lua"
	"github.com/colinyl/servicebus/rpc"
	l "github.com/yuin/gopher-lua"
)

var AppLuaEngine *lua.LuaPool

func RPCRequest(L *l.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]l.LGFunction{
	"request": request,
}

func request(L *l.LState) int {
	name := L.ToString(1)
	input := L.ToString(2)
	result, err := rpc.RCServerPool.Request(name, input)
	L.Push(l.LString(result))
	if err != nil {
		L.Push(l.LString(err.Error()))
	} else {
		L.Push(l.LNil)
	}
	return 2
}

func init() {
	rpcFunc := lua.Luafunc{Name: "rpc", Function: RPCRequest}
	AppLuaEngine = lua.NewLuaPool(rpcFunc)
}
