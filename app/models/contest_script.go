package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/yuin/gopher-lua"
)

var lCodes = map[int32]string{}
var lStates = map[int32]*lua.LState{}

func (c *Contest) ResetLuaState() {
	L := lStates[c.Id]
	if L != nil {
		L.Close()
	}

	L = lua.NewState()
	lStates[c.Id] = L

	L.SetGlobal("get_handle", L.NewFunction(luaGetHandle))
	L.SetGlobal("get_id", L.NewFunction(luaGetId))
	L.SetGlobal("create_match", L.NewFunction(luaCreateMatch))

	err := L.DoString(c.Script)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (c *Contest) LuaState() *lua.LState {
	if lCodes[c.Id] != c.Script {
		lCodes[c.Id] = c.Script
		c.ResetLuaState()
	}
	return lStates[c.Id]
}

// Finds a user's handle according to their ID
// If the ID is not found or an error happens, an empty string is returned
func luaGetHandle(L *lua.LState) int {
	// Argument
	uid := L.ToInt(1)

	u := User{Id: int32(uid)}
	if err := u.ReadById(); err != nil {
		L.Push(lua.LString(""))
		return 1
	}

	L.Push(lua.LString(u.Handle))
	return 1
}

// Finds a user's ID according to their handle
// If the handle is not found or an error happens, 0 is returned
func luaGetId(L *lua.LState) int {
	// Argument
	handle := L.ToString(1)

	u := User{Handle: handle}
	if err := u.ReadByHandle(); err != nil {
		L.Push(lua.LNumber(0))
		return 1
	}

	L.Push(lua.LNumber(u.Id))
	return 1
}

func luaCreateMatch(L *lua.LState) int {
	// Arguments
	argc := L.GetTop()
	ss := []Submission{}
	for i := 1; i <= argc; i++ {
		sid := L.ToInt(i)
		ss = append(ss, Submission{Id: int32(sid)})
	}

	// Create match
	for i, s := range ss {
		fmt.Printf("[%d %d]\n", i, s.Id)
	}

	return 0
}

func (c *Contest) ExecuteScriptOnTimer() error {
	L := c.LuaState()

	// Find `on_timer` global function
	val := L.GetGlobal("on_timer")
	if val.Type() != lua.LTFunction {
		return errors.New("Lua global `on_timer` should be a function")
	}

	// Retrieve all contestants and make a Lua table
	ps, err := c.AllParticipations()
	if err != nil {
		return err
	}
	t := &lua.LTable{}
	for _, p := range ps {
		if p.Delegate != -1 {
			t.Append(lua.LNumber(p.Rel.User.Id))
		}
	}

	// Call Lua function
	err = L.CallByParam(lua.P{
		Fn:      val,
		NRet:    0,
		Protect: true,
	}, t)
	if err != nil {
		return err
	}

	return nil
}

func timerForAllContests() {
	for {
		// println("Script timer")
		time.Sleep(2 * time.Second)
		cs, err := ContestReadAll()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		for _, c := range cs {
			if c.IsRunning() {
				if err := c.ExecuteScriptOnTimer(); err != nil {
					fmt.Println(err.Error())
					continue
				}
			}
		}
	}
}

func init() {
	go timerForAllContests()
}
