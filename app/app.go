package app

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"bitbucket.org/syb-devs/goth/database"
	"bitbucket.org/syb-devs/goth/encoding/json"
	"bitbucket.org/syb-devs/goth/kv"
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
	Muxer
	modules map[string]Module
	mws     map[string]MiddlewareChain
}

// NewApp instances and returns an App with the given name
func NewApp(name string) *App {
	return &App{
		name:    name,
		modules: modules,
		mws:     make(map[string]MiddlewareChain),
	}
}

// Name returns the App name
func (a *App) Name() string {
	return a.name
}

// AddChain adds a MiddlewareChain to be used when registering handlers
func (a *App) AddChain(chain MiddlewareChain, name string) {
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

// SetMuxer sets a Muxer for the App object
func (a *App) SetMuxer(r Muxer) {
	a.Lock()
	defer a.Unlock()
	a.Muxer = r
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

	fmt.Printf("### %s is starting...\n", a.name)
	a.bootstrap()
	fmt.Printf("listening on port 8080...\n")
	log.Fatal(http.ListenAndServe(":8080", a))
}

// Close closes the database connection of this App instance (Copy/Close pattern)
func (a *App) Close() {
	a.DB.Close()
}

// RegisterResource registers a new Resource in the App
func (a *App) RegisterResource(res database.Resource, conf database.ResourceConfig) {
	if conf.Name == "" {
		panic("trying to register resource without name")
	}
	// TODO(zareone) define a RegisterResource hook / callback in the module interface
	// and perform this tasks from the MongoDB and HTTP modules??
	a.registerResourceDB(res, conf)
}

func (a *App) registerResourceDB(res database.Resource, conf database.ResourceConfig) error {
	if conf.ColName == "" {
		panic(fmt.Sprintf("invalid collection name for resource %s", conf.Name))
	}
	var deletedCol string
	if conf.ArchiveOnDelete {
		deletedCol = "archived_" + conf.ColName
	}
	return a.DB.ResourceMap.RegisterResource(res, conf.ColName, deletedCol)
}

func (a *App) bootstrap() {
	for level := 0; level <= 10; level++ {
		for name, mod := range a.modules {
			fmt.Printf("bootstrapping module %s at level %d\n", name, level)
			err := mod.Bootstrap(a, level)
			if err != nil {
				panic(fmt.Sprintf("bootstrap error: module: %s, level: %d", name, level))
			}
		}
	}
}
