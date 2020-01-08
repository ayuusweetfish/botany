package models

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/yuin/gopher-lua"
)

type ErrLuaType struct {
	Inner   error
	Message string
}

var ErrLua = errors.New("")

func (e ErrLuaType) Unwrap() error {
	return ErrLua
}

func (e ErrLuaType) Error() string {
	return e.Message
}

var lCodes = map[int32]string{}
var lStates = map[int32]*lua.LState{}
var lCIDs = map[*lua.LState]int32{}
var lLogs = map[*lua.LState]*strings.Builder{}
var lLogTail = map[int32]*tailer{}

const tailNumLines = 50

type tailer struct {
	Lines [tailNumLines]string
	Count uint32
	Ptr   uint16
}

func (t *tailer) Append(s string) {
	t.Lines[t.Ptr] = s
	t.Ptr++
	t.Count++
	if t.Ptr == tailNumLines {
		t.Ptr = 0
	}
}

func (t *tailer) Join() string {
	b := strings.Builder{}
	if t.Count < tailNumLines {
		for i := uint32(0); i < t.Count; i++ {
			b.WriteString(t.Lines[i])
			b.WriteRune('\n')
		}
	} else {
		for i := uint16(0); i < tailNumLines; i++ {
			b.WriteString(t.Lines[(i+t.Ptr)%tailNumLines])
			b.WriteRune('\n')
		}
	}
	return b.String()
}

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
		log.Println(err.Error())
	}

	flushLog(c.Id)
}

func (c *Contest) LuaState() *lua.LState {
	if lCodes[c.Id] != c.Script || lStates[c.Id] == nil {
		lCodes[c.Id] = c.Script
		c.ResetLuaState()
	}
	return lStates[c.Id]
}

func writeTimestamp(b *strings.Builder) {
	b.WriteRune('[')
	b.WriteString(time.Now().Format("2006-01-02 15:04:05 -0700 MST"))
	b.WriteRune(']')
	b.WriteRune('\t')
}

func flushLog(cid int32) {
	// Retrieve strings.Builder
	builder := lLogs[lStates[cid]]
	if builder == nil {
		lLogs[lStates[cid]] = &strings.Builder{}
		return
	}
	s := builder.String()

	// Update database
	c := Contest{Id: cid}
	if err := c.AppendScriptLog(s); err != nil {
		log.Println(err.Error())
		return
	}

	// Update in-memory cache
	// TODO: Cache needs to be populated on startup
	t := lLogTail[cid]
	if t == nil {
		t = &tailer{Count: 0, Ptr: 0}
		lLogTail[cid] = t
	}
	if s != "" && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	if s != "" {
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			t.Append(line)
		}
	}

	builder.Reset()
}

func (c *Contest) TailLog() string {
	t := lLogTail[c.Id]
	if t == nil {
		return ""
	}
	return t.Join()
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
			log.Println(err.Error())
			return 0
		}
		p := ContestParticipation{User: uid, Contest: cid}
		if err := p.Read(); err != nil {
			log.Println(err.Error())
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
		Report:  "Pending",
	}
	m.Rel.Parties = ss
	if err := m.Create(); err != nil {
		log.Println(err.Error())
		return 0
	}

	// Get judge ID and then add to judge queue
	c := Contest{Id: cid}
	if err := c.Read(); err != nil {
		log.Println(err.Error())
		return 0
	}
	if err := m.SendToQueue(c.Judge); err != nil {
		log.Println(err.Error())
		return 0
	}

	builder.WriteString(fmt.Sprintf(" - #%d", m.Id))
	return 0
}

func (c *Contest) ExecuteScriptFunction(fnName string, args ...lua.LValue) error {
	L := c.LuaState()

	// Find global function by name
	fn := L.GetGlobal(fnName)
	if fn.Type() != lua.LTFunction {
		return ErrLuaType{Message: "Lua global `" + fnName + "` should be a function"}
	}

	// Retrieve all contestants and make a Lua table
	ps, err := c.AllParticipationsWithDelegate()
	if err != nil {
		return err
	}
	t := &lua.LTable{}
	for _, p := range ps {
		t.Append(lua.LNumber(p.Rel.User.Id))
	}

	// Call Lua function
	err = L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    0,
		Protect: true,
	}, append([]lua.LValue{t}, args...)...)
	if err != nil {
		return err
	}

	flushLog(c.Id)
	return nil
}

func (c *Contest) ExecuteScriptOnSubmission(from int32) error {
	return c.ExecuteScriptFunction("on_submission", lua.LNumber(from))
}

func (c *Contest) ExecuteScriptOnTimer() error {
	return c.ExecuteScriptFunction("on_timer")
}

func (c *Contest) ExecuteScriptOnManual(arg string) error {
	return c.ExecuteScriptFunction("on_manual", lua.LString(arg))
}

func timerForAllContests() {
	for {
		time.Sleep(2 * time.Second)
		cs, err := ContestReadAll()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		for _, c := range cs {
			if c.IsRunning() {
				if err := c.ExecuteScriptOnTimer(); err != nil {
					log.Println(err.Error())
					continue
				}
			}
		}
	}
}

func init() {
	go timerForAllContests()
}
