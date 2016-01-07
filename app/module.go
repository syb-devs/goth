package app

import (
	"fmt"
)

var modules = make(map[string]Module)

// Module is used to extend the base App with new functionalities
type Module interface {
	Name() string
	Bootstrap(*App, int) error
}

// RegisterModule is used by external modulse to register themselves on the App
func RegisterModule(module Module) {
	name := module.Name()
	if _, exists := modules[name]; exists {
		panic(fmt.Sprintf("module %s already registered", name))
	}
	modules[name] = module
}

// NewBaseModule allocates and returns a BaseModule
func NewBaseModule(name string) *BaseModule {
	return &BaseModule{
		name: name,
	}
}

// BaseModule is used as a base to implement the module interface with some defaults
type BaseModule struct {
	name string
}

func (m *BaseModule) mustImplement(method string) {
	panic(fmt.Sprintf("module %s must implement %s", m.name, method))
}

// Name default implementation
func (m *BaseModule) Name() string { return m.name }

// Bootstrap default implementation
func (m *BaseModule) Bootstrap(*App, int) error { m.mustImplement("Bootstrap"); return nil }
