package rest

import (
	"fmt"

	"bitbucket.org/syb-devs/goth/app"
	"bitbucket.org/syb-devs/goth/database"
)

const (
	resourceIDParam = "resource_id"
)

// ResourceHandler interface represents the HTTP interface for CRUD operations
// than can be applied to a Resource
type ResourceHandler interface {
	Create(ctx *app.Context) error
	Retrieve(ctx *app.Context) error
	Update(ctx *app.Context) error
	Delete(ctx *app.Context) error
	List(ctx *app.Context) error
}

// ResourceConfig has the needed info to Register a resource in the App
type ResourceConfig struct {
	Name    string
	URLName string
	Handler ResourceHandler
}

// Register registers a resource for setting up automatic REST CRUD handlers in the App
func Register(a *app.App, conf ResourceConfig) {
	pName := conf.URLName
	if pName == "" {
		panic(fmt.Sprintf("empty path name for resource %s", conf.Name))
	}
	URL := fmt.Sprintf("/%s", pName)
	URLWithID := fmt.Sprintf("/%s/:%s", pName, resourceIDParam)
	rh := conf.Handler

	// Register CRUD routes for Resource
	a.Handle("POST", URL, app.HandlerFunc(rh.Create))
	a.Handle("GET", URLWithID, app.HandlerFunc(rh.Retrieve))
	a.Handle("PUT", URLWithID, app.HandlerFunc(rh.Update))
	a.Handle("DELETE", URLWithID, app.HandlerFunc(rh.Delete))
	a.Handle("GET", URL, app.HandlerFunc(rh.List))
}

// DefResourceHandler is the default implementation for ResourceHandler interface
type DefResourceHandler struct {
	constructor database.ResourceConstructor
	database.ResourceValidator
}

// NewDefResourceHandler allocates and returns a DefResourceHandler
func NewDefResourceHandler(rc database.ResourceConstructor) *DefResourceHandler {
	if rc == nil {
		panic("valid resource constructor needed")
	}

	return &DefResourceHandler{
		constructor:       rc,
		ResourceValidator: &database.DummyResourceValidator{},
	}
}

// Create decodes a resource from the Request, validates it and stores it in the database
func (h *DefResourceHandler) Create(ctx *app.Context) error {
	res := h.constructor.New()
	err := ctx.Decode(res)
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
	return ctx.Encode(res)
}

// Retrieve fetches a resource from the database and encodes it to the ResponseWriter
func (h *DefResourceHandler) Retrieve(ctx *app.Context) error {
	ID := ctx.URLParams.ByName(resourceIDParam)
	res := h.constructor.New()
	err := ctx.App.DB.Get(ID, res)
	if err != nil {
		return err
	}
	return ctx.Encode(res)
}

// Update decodes a resource from the Request, validates it and updates it in the database
func (h *DefResourceHandler) Update(ctx *app.Context) error {
	res := h.constructor.New()
	ID := ctx.URLParams.ByName(resourceIDParam)
	err := ctx.App.DB.Get(ID, res)
	if err != nil {
		return err
	}
	err = ctx.Decode(res)
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
	return ctx.Encode(res)
}

// Delete deletes a Resource from the database
func (h *DefResourceHandler) Delete(ctx *app.Context) error {
	res := h.constructor.New()
	ID := ctx.URLParams.ByName(resourceIDParam)
	err := ctx.App.DB.Get(ID, res)
	if err != nil {
		return err
	}
	return ctx.App.DB.Delete(res)
}

// List retrieves a list of Resources from the database, and encodes it to the ResponseWriter
func (h *DefResourceHandler) List(ctx *app.Context) error {
	list := h.constructor.NewList()
	query := database.NewQ(nil)
	err := ctx.App.DB.FindMany(list, query)
	if err != nil {
		return err
	}
	return ctx.Encode(list)
}
