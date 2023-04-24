package scripts

import lua "github.com/yuin/gopher-lua"

var luaInstanceConvert *lua.LState

func LoadFromFileConvert(luaFile string) {
	luaInstanceConvert = lua.NewState()
	err := luaInstanceConvert.DoFile(luaFile)
	if err != nil {
		panic(err)
	}
}

func Convert(table, ddl string) string {
	if luaInstanceConvert == nil {
		return ddl
	}

	f := luaInstanceConvert.GetGlobal("convert")
	luaInstanceConvert.Push(f)
	luaInstanceConvert.Push(lua.LString(table)) // table name
	luaInstanceConvert.Push(lua.LString(ddl))   // table name

	luaInstanceConvert.Call(2, 1)

	newDDL := string(luaInstanceConvert.Get(1).(lua.LString))
	luaInstanceConvert.Pop(1)
	return newDDL
}
