package rest

import (
	"fmt"
	"net/http"

	"github.com/syb-devs/goth/app"
	"github.com/syb-devs/goth/database"
)

// ResourceHandler interface represents the HTTP interface for CRUD operations
// than can be applied to a Resource
type ResourceHandler interface {
	Create(w http.ResponseWriter, r *http.Request, ctx *app.Context) error
	Retrieve(w http.ResponseWriter, r *http.Request, ctx *app.Context) error
	Update(w http.ResponseWriter, r *http.Request, ctx *app.Context) error
	Delete(w http.ResponseWriter, r *http.Request, ctx *app.Context) error
	List(w http.ResponseWriter, r *http.Request, ctx *app.Context) error
}

// ResourceConfig has the needed info to Register a resource in the App
type ResourceConfig struct {
	Name    string
	URLName string
	Handler ResourceHandler
}

func RegisterResource(a *app.App, conf ResourceConfig) error {
	pName := conf.URLName
	if pName == "" {
		panic(fmt.Sprintf("empty path name for resource %s", conf.Name))
	}
	URL := fmt.Sprintf("/%s", pName)
	URLWithID := fmt.Sprintf("/%s/:id", pName)
	rh := conf.Handler

	// Register CRUD routes for Resource
	a.Handle("POST", URL, app.HandlerFunc(rh.Create))
	a.Handle("GET", URLWithID, app.HandlerFunc(rh.Retrieve))
	a.Handle("PUT", URLWithID, app.HandlerFunc(rh.Update))
	a.Handle("DELETE", URLWithID, app.HandlerFunc(rh.Delete))
	a.Handle("GET", URL, app.HandlerFunc(rh.List))

	if rp, ok := rh.(app.RouteProvider); ok {
		// Register any extra routes defined by the ResourceHandler
		rp.RegisterRoutes(a)
	}
	return nil
}

// DefResourceHandler is the default implementation for ResourceHandler interface
type DefResourceHandler struct {
	constructor database.ResourceConstructor
	database.ResourceCodec
	database.ResourceValidator
}

// NewDefResourceHandler allocates and returns a DefResourceHandler
func NewDefResourceHandler(rc database.ResourceConstructor) *DefResourceHandler {
	if rc == nil {
		panic("valid resource constructor needed")
	}

	return &DefResourceHandler{
		constructor:       rc,
		ResourceCodec:     &database.JSONResourceCodec{},
		ResourceValidator: &database.DummyResourceValidator{},
	}
}

// Create decodes a resource from the Request, validates it and stores it in the database
func (h *DefResourceHandler) Create(w http.ResponseWriter, r *http.Request, ctx *app.Context) error {
	res := h.constructor.New()
	err := h.Decode(r.Body, res)
	if err != nil {
		return err
	}
	err = h.Validate(database.ResourceActionCreate, res)
	if err != nil {
		return err
	}
	err = ctx.App.DB.Insert(res)
	if err != nil {
		return err
	}
	return h.Encode(w, res)
}

// Retrieve fetches a resource from the database and encodes it to the ResponseWriter
func (h *DefResourceHandler) Retrieve(w http.ResponseWriter, r *http.Request, ctx *app.Context) error {
	ID := ctx.URLParams.ByName("id")
	res := h.constructor.New()
	err := ctx.App.DB.Get(ID, res)
	if err != nil {
		return err
	}
	return h.Encode(w, res)
}

// Update decodes a resource from the Request, validates it and updates it in the database
func (h *DefResourceHandler) Update(w http.ResponseWriter, r *http.Request, ctx *app.Context) error {
	res := h.constructor.New()
	ID := ctx.URLParams.ByName("id")
	err := ctx.App.DB.Get(ID, res)
	if err != nil {
		return err
	}
	err = h.Decode(r.Body, res)
	if err != nil {
		return err
	}
	err = h.Validate(database.ResourceActionUpdate, res)
	if err != nil {
		return err
	}
	err = ctx.App.DB.Update(res)
	if err != nil {
		return err
	}
	return h.Encode(w, res)
}

// Delete deletes a Resource from the database
func (h *DefResourceHandler) Delete(w http.ResponseWriter, r *http.Request, ctx *app.Context) error {
	res := h.constructor.New()
	ID := ctx.URLParams.ByName("id")
	err := ctx.App.DB.Get(ID, res)
	if err != nil {
		return err
	}
	err = ctx.App.DB.Delete(res)
	if err != nil {
		return err
	}
	return nil
}

// List retrieves a list of Resources from the database, and encodes it to the ResponseWriter
func (h *DefResourceHandler) List(w http.ResponseWriter, r *http.Request, ctx *app.Context) error {
	list := h.constructor.NewList()
	query := database.NewQ(nil)
	fmt.Println("list!")
	err := ctx.App.DB.FindMany(list, query)
	if err != nil {
		return err
	}
	return h.Encode(w, list)
}
