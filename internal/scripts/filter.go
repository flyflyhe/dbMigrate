package scripts

import (
	lua "github.com/yuin/gopher-lua"
)

const (
	Allow    = 0
	Disallow = 1
	Error    = 2
)

var luaInstance *lua.LState

func LoadFromFile(luaFile string) {
	luaInstance = lua.NewState()
	err := luaInstance.DoFile(luaFile)
	if err != nil {
		panic(err)
	}
}

func Filter(table string) int {
	if luaInstance == nil {
		return Allow
	}

	f := luaInstance.GetGlobal("filter")
	luaInstance.Push(f)
	luaInstance.Push(lua.LString(table)) // table name

	luaInstance.Call(1, 1)

	code := int(luaInstance.Get(1).(lua.LNumber))
	luaInstance.Pop(1)
	return code
}
