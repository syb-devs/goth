package app

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"bitbucket.org/syb-devs/goth/database"
	"bitbucket.org/syb-devs/goth/encoding/json"
	"bitbucket.org/syb-devs/goth/kv"
	"bitbucket.org/syb-devs/goth/log"
)

// App represents the main application
type App struct {
	sync.Mutex
	name string
	DB   struct {
		database.Connection
		database.Repository
		*database.ResourceMap
	}
	Log     log.Logger
	muxer   Muxer
	handler http.Handler
	modules map[string]Module
	mws     map[string]MiddlewareChain
}

// NewApp instances and returns an App with the given name
func NewApp(name string) *App {
	return &App{
		name:    name,
		modules: modules,
		mws:     make(map[string]MiddlewareChain),
		Log:     log.New(os.Stderr),
	}
}

// Name returns the App name
func (a *App) Name() string {
	return a.name
}

// AddChain adds a MiddlewareChain to be used when registering handlers
func (a *App) AddChain(chain MiddlewareChain, name string) {
	a.Lock()
	defer a.Unlock()
	a.mws[name] = chain
}

// WrapHandler wraps the given Handler with a MiddlewareChain
func (a *App) WrapHandler(h Handler, chainName string) Handler {
	chain, ok := a.mws[chainName]
	if !ok {
		panic(fmt.Errorf("no middleware chain registered with name %s\n", chainName))
	}
	return chain.Finally(h)
}

// WrapHandlerFunc wraps the given HandlerFunc with a MiddlewareChain
func (a *App) WrapHandlerFunc(h HandlerFunc, chainName string) Handler {
	return a.WrapHandler(h, chainName)
}

// Handle registers a handler for a given method / path combination
func (a *App) Handle(method, path string, h Handler) {
	log.Debugf("registering route %s %s", method, path)
	a.muxer.Handle(method, path, h)
}

// SetMuxer sets a Muxer for the App object
func (a *App) SetMuxer(m Muxer) {
	a.Lock()
	defer a.Unlock()
	a.muxer = m
}

// SetHandler sets a Handler for the App object
func (a *App) SetHandler(h http.Handler) {
	a.Lock()
	defer a.Unlock()
	a.handler = h
}

// NewContextHTTP creates a new context for the given HTTP request
func (a *App) NewContextHTTP(w http.ResponseWriter, r *http.Request) *Context {
	//TODO(zareone): init DB and Codec
	return &Context{
		App:            a,
		Request:        r,
		ResponseWriter: w,
		Store:          kv.New(),
		Codec:          json.Codec{},
	}
}

// Run starts the app server
func (a *App) Run() {
	defer a.Close()

	a.Log.Infof("### %s is starting...\n", a.name)
	a.bootstrap()
	a.Log.Infof("listening on port 8080...\n")
	err := http.ListenAndServe(":8080", a.handler)
	a.Log.Error(err)
	os.Exit(1)
}

// Close closes the database connection of this App instance (Copy/Close pattern)
func (a *App) Close() {
	a.DB.Close()
}

func (a *App) bootstrap() {
	for level := 0; level <= 10; level++ {
		for name, mod := range a.modules {
			a.Log.Debugf("bootstrapping module %s at level %d\n", name, level)
			err := mod.Bootstrap(a, level)
			if err != nil {
				panic(fmt.Sprintf("bootstrap error: module: %s, level: %d", name, level))
			}
		}
	}
}
