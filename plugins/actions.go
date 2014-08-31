package plugins

import (
	"sort"
)

type ActionFn func(args map[string]interface{}) map[string]interface{}
type Action struct {
	Fn       ActionFn
	Priority int
	ArgCount int
}

type Actions []Action

var _actions map[string][]Action = make(map[string][]Action)

func AddAction(name string, fn ActionFn, priority, args int) {
	var act Action
	if _actions[name] == nil {
		_actions[name] = make([]Action, 0)
	}
	act.Fn = fn
	act.Priority = priority
	act.ArgCount = args
	_actions[name] = append(_actions[name], act)

	sort.Sort(sort.Reverse(Actions(_actions[name])))
}

func HasAction(name string) bool {
	return false
}

func DoActions(name string, arg map[string]interface{}) {
	if _actions[name] == nil {
		return
	}

	for _, action := range _actions[name] {
		action.Fn(arg)
	}
}

// 排序
func (a Actions) Len() int {
	return len(a)
}

func (a Actions) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Actions) Less(i, j int) bool {
	return a[i].Priority < a[j].Priority
}
