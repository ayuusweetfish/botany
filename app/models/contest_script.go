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
		t.Append(lua.LNumber(p.Rel.User.Id))
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
		println("Script timer")
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
