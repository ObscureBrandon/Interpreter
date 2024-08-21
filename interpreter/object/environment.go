package object

import "fmt"

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Set(name string, val Object, is_declaration bool) Object {
	_, ok := e.store[name]

	if is_declaration && ok {
		return &Error{Message: fmt.Sprintf("%q declared in for loop initilization", name)}
	}

	if is_declaration {
		e.store[name] = val
		return val
	}

	if !is_declaration {
		_, ok = e.store[name]
		if ok {
			e.store[name] = val
			return val
		}

		if e.outer != nil {
			e.outer.store[name] = val
			return val
		}
	}

	return val
}
