package app

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/syb-devs/goth/database"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/justinas/alice"
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
	Router
	mws     []middlewareOpts
	gmws    []Middleware
	modules map[string]Module
}

// NewApp instances and returns an App with the given name
func NewApp(name string) *App {
	return &App{
		name:    name,
		modules: modules,
	}
}

// Name returns the App name
func (a *App) Name() string {
	return a.name
}

// SetRouter sets a router for the App object
func (a *App) SetRouter(r Router) {
	a.Lock()
	defer a.Unlock()
	a.Router = r
}

// Use registers a middleware with an alias name to be used in the app
func (a *App) Use(mw Middleware, alias string) {
	a.Lock()
	defer a.Unlock()

	a.mws = append(a.mws, middlewareOpts{alias, mw})
}

// UseGlobal registers a  global middleware that will wrap the App Router
func (a *App) UseGlobal(mw Middleware) {
	a.Lock()
	defer a.Unlock()

	a.gmws = append(a.gmws, mw)
}

func (a *App) middlewareToHandler(h Handler) Handler {
	if len(a.mws) == 0 {
		return addContext(h, a)
	}
	var mws = make([]alice.Constructor, 0, len(a.mws))

	skipper, isSkipper := h.(MWSkipper)
	for _, mwo := range a.mws {
		if isSkipper && skipper.SkipMiddleware(mwo.alias) {
			continue
		}
		mws = append(mws, alice.Constructor(mwo.middleware))
	}
	chain := alice.New(mws...)
	return adaptHandler(chain.Then(wrapHandler(h, a)))
}

func (a *App) middlewareToRouter(r Router) http.Handler {
	if len(a.gmws) == 0 {
		return r
	}
	var cs = make([]alice.Constructor, 0, len(a.gmws))
	for _, mw := range a.gmws {
		cs = append(cs, alice.Constructor(mw))
	}
	return alice.New(cs...).Then(r)
}

// Handle registers an HTTP handler to a given verb / path combination
func (a *App) Handle(verb, path string, h Handler) {
	a.Lock()
	defer a.Unlock()

	fmt.Printf("registering route %s %s\n", verb, path)
	if a.Router == nil {
		panic("a router must be set to the app")
	}
	a.Router.Handle(verb, path, a.middlewareToHandler(h))
}

func (a *App) newContext(w http.ResponseWriter, r *http.Request) *Context {
	ctxApp := a.Copy()

	var uid string
	if token, ok := context.Get(r, "user").(*jwt.Token); ok {
		uid = token.Claims["user_id"].(string)
	}
	return &Context{
		App:       ctxApp,
		URLParams: context.Get(r, "url_params").(URLParams),
		Conn:      ctxApp.DB.Connection,
		UserID:    uid,
		Request:   r,
	}
}

// Run starts the app server
func (a *App) Run() {
	fmt.Printf("### %s is starting...\n", a.name)
	a.bootstrap()
	a.Handle("GET", "/", SkipMiddlewareFunc(rootHandler, "jwt"))
	fmt.Printf("listening on port 8080...\n")
	log.Fatal(http.ListenAndServe(":8080", a.middlewareToRouter(a.Router)))
}

// Copy creates a shallow copy of the App, with a copied database connection
func (a *App) Copy() *App {
	appCopy := *a
	appCopy.DB.Connection = a.DB.Copy()
	return &appCopy
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

// Middleware is a function that wraps an HTTP handler for doing some work before or after it
type Middleware func(http.Handler) http.Handler

// middlewareOpts is used to register aliased middlewares
type middlewareOpts struct {
	alias      string
	middleware Middleware
}

// MWSkipper interface allows handlers to skip ms
type MWSkipper interface {
	SkipMiddleware(string) bool
}

// MWSkipperHandler is a Handler that has been configured to skip some ms
type MWSkipperHandler struct {
	Handler
	mws []string
}

// SkipMiddleware returns a MWSkipperHandler to skip the specified middleware(s)
func SkipMiddleware(h Handler, aliases ...string) *MWSkipperHandler {
	return &MWSkipperHandler{
		Handler: h,
		mws:     aliases,
	}
}

// SkipMiddlewareFunc returns a MWSkipperHandler to skip the specified middleware(s)
func SkipMiddlewareFunc(h HandlerFunc, aliases ...string) *MWSkipperHandler {
	return &MWSkipperHandler{
		Handler: h,
		mws:     aliases,
	}
}

// SkipMiddleware checks if the Handler should skip the given middleware
func (sm *MWSkipperHandler) SkipMiddleware(alias string) bool {
	for _, mw := range sm.mws {
		if alias == mw {
			return true
		}
	}
	return false
}
