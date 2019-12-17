package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/yuin/gopher-lua"
)

var lCodes = map[int32]string{}
var lStates = map[int32]*lua.LState{}
var lCIDs = map[*lua.LState]int32{}
var lLogs = map[*lua.LState]*strings.Builder{}

func (c *Contest) ResetLuaState() {
	L := lStates[c.Id]
	if L != nil {
		L.Close()
	}

	L = lua.NewState()
	lStates[c.Id] = L
	lCIDs[L] = c.Id
	flushLog(c.Id)

	L.SetGlobal("print", L.NewFunction(luaBasePrintRedirect))

	L.SetGlobal("get_handle", L.NewFunction(luaGetHandle))
	L.SetGlobal("get_id", L.NewFunction(luaGetId))
	L.SetGlobal("create_match", L.NewFunction(luaCreateMatch))

	err := L.DoString(c.Script)
	if err != nil {
		fmt.Println(err.Error())
	}

	flushLog(c.Id)
}

func (c *Contest) LuaState() *lua.LState {
	if lCodes[c.Id] != c.Script {
		lCodes[c.Id] = c.Script
		c.ResetLuaState()
	}
	return lStates[c.Id]
}

func writeTimestamp(builder *strings.Builder) {
	builder.WriteRune('[')
	builder.WriteString(time.Now().Format("2006-01-02 15:04:05 -0700 MST"))
	builder.WriteRune(']')
	builder.WriteRune('\t')
}

func flushLog(cid int32) {
	builder := lLogs[lStates[cid]]
	if builder == nil {
		lLogs[lStates[cid]] = &strings.Builder{}
		return
	}

	c := Contest{Id: cid}
	if err := c.AppendScriptLog(builder.String()); err != nil {
		fmt.Println(err.Error())
		return
	}

	builder.Reset()
}

func luaBasePrintRedirect(L *lua.LState) int {
	builder := lLogs[L]
	writeTimestamp(builder)

	top := L.GetTop()
	for i := 1; i <= top; i++ {
		builder.WriteString(L.ToStringMeta(L.Get(i)).String())
		if i != top {
			builder.WriteRune('\t')
		}
	}
	builder.WriteRune('\n')
	return 0
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
	cid := lCIDs[L]

	// Arguments
	argc := L.GetTop()
	ss := []Submission{}

	builder := lLogs[L]
	writeTimestamp(builder)
	builder.WriteString("<Match created among")
	defer builder.WriteString(">\n")

	for i := 1; i <= argc; i++ {
		uid := int32(L.ToInt(i))
		u := User{Id: uid}
		if err := u.ReadById(); err != nil {
			fmt.Println(err.Error())
			return 0
		}
		p := ContestParticipation{User: uid, Contest: cid}
		if err := p.Read(); err != nil {
			fmt.Println(err.Error())
			return 0
		}
		sid := p.Delegate
		if i > 1 {
			builder.WriteRune(',')
		}
		builder.WriteRune(' ')
		builder.WriteString(u.Handle)
		ss = append(ss, Submission{Id: sid})
	}

	// Create match
	m := Match{
		Contest: cid,
		Report:  "{\"winner\": \"In queue\"}",
	}
	m.Rel.Parties = ss
	if err := m.Create(); err != nil {
		fmt.Println(err.Error())
		return 0
	}
	if err := m.SendToQueue(); err != nil {
		fmt.Println(err.Error())
		return 0
	}

	builder.WriteString(fmt.Sprintf(" - #%d", m.Id))
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

	flushLog(c.Id)
	return nil
}

func timerForAllContests() {
	for {
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
