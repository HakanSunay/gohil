package env

import "github.com/HakanSunay/gohil/object"

type Environment struct {
	store map[string]object.Object
}

func NewEnvironment() *Environment {
	s := make(map[string]object.Object)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (object.Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val object.Object) object.Object {
	e.store[name] = val
	return val
}
