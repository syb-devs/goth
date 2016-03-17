package rest

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"bitbucket.org/syb-devs/goth/app"
	"bitbucket.org/syb-devs/goth/database"
)

const (
	resourceIDParam = "resource_id"
)

// ErrNotAllowed
var ErrNotAllowed = errors.New("operation not allowed")

// CRUDHandler interface represents the HTTP interface for CRUD operations
// than can be applied to a Resource
type CRUDHandler interface {
	Create(ctx *app.Context) error
	Retrieve(ctx *app.Context) error
	Update(ctx *app.Context) error
	Delete(ctx *app.Context) error
	List(ctx *app.Context) error
}

// Register registers a resource for setting up automatic REST CRUD handlers in the App
func Register(a *app.App, res database.Resource, name string) {
	RegisterCRUD(a, New(res), name)
}

// RegisterWithOptions registers a resource for setting up automatic REST CRUD handlers in the App, using the given Options
func RegisterWithOptions(a *app.App, res database.Resource, opts Options, name string) {
	RegisterCRUD(a, NewWithOptions(res, opts), name)
}

// RegisterCRUD registers a CRUDHandler
func RegisterCRUD(a *app.App, crud CRUDHandler, name string) {
	URL := fmt.Sprintf("/%s", name)
	URLWithID := fmt.Sprintf("/%s/:%s", name, resourceIDParam)

	// Register CRUD routes for Resource
	a.Handle("POST", URL, a.WrapHandlerFunc(crud.Create, "main"))
	a.Handle("GET", URLWithID, a.WrapHandlerFunc(crud.Retrieve, "main"))
	a.Handle("PUT", URLWithID, a.WrapHandlerFunc(crud.Update, "main"))
	a.Handle("DELETE", URLWithID, a.WrapHandlerFunc(crud.Delete, "main"))
	a.Handle("GET", URL, a.WrapHandlerFunc(crud.List, "main"))
}

// BaseCRUD is the default implementation for ResourceHandler interface
type BaseCRUD struct {
	resourceLimit int
	resourceType  reflect.Type
	ownerField    string
	validator     database.ResourceValidator
	accessChecker app.ResourceAccessChecker
}

// Options are used to setup the CRUD
type Options struct {
	ResourceLimit int
	OwnerDBField  string
	AccessChecker app.ResourceAccessChecker
}

// New allocates and returns a BaseCRUD
func New(resource database.Resource) *BaseCRUD {
	return &BaseCRUD{
		resourceLimit: 100,
		resourceType:  reflect.TypeOf(resource),
		ownerField:    "",
		validator:     &database.DummyResourceValidator{},
		accessChecker: app.AccessChecker{},
	}
}

// NewWithOptions allocates and returns a BaseCRUD using the given Options
func NewWithOptions(resource database.Resource, opts Options) *BaseCRUD {
	crud := New(resource)
	if opts.AccessChecker != nil {
		crud.accessChecker = opts.AccessChecker
	}
	if opts.ResourceLimit != 0 {
		crud.resourceLimit = opts.ResourceLimit
	}
	if opts.OwnerDBField != "" {
		crud.ownerField = opts.OwnerDBField
	}
	return crud
}

// NewResource allocates a new object of the type of the resource
func (h *BaseCRUD) NewResource(ctx *app.Context) database.Resource {
	return ctx.App.DB.CreateResource(h.resourceType).(database.Resource)
}

// NewResourceList allocates a new slice of objects of the type of the resource
func (h *BaseCRUD) NewResourceList(ctx *app.Context) database.ResourceList {
	return ctx.App.DB.CreateResourceList(h.resourceType)
}

// Create decodes a resource from the Request, validates it and stores it in the database
func (h *BaseCRUD) Create(ctx *app.Context) error {
	res := h.NewResource(ctx)
	err := ctx.Decode(res)
	if err != nil {
		return err
	}
	err = h.validator.Validate(database.ResourceActionCreate, res)
	if err != nil {
		return err
	}
	res.SetOwnerID(ctx.User.GetID())
	err = ctx.App.DB.Insert(res)
	if err != nil {
		return err
	}
	return ctx.Encode(res)
}

// Retrieve fetches a resource from the database and encodes it to the ResponseWriter
func (h *BaseCRUD) Retrieve(ctx *app.Context) error {
	ID := ctx.URLParams.ByName(resourceIDParam)
	res := h.NewResource(ctx)
	err := ctx.App.DB.Get(ID, res)
	if err != nil {
		return err
	}
	err = h.checkAccess(ctx, res, "view")
	if err != nil {
		return err
	}
	err = expandResource(ctx, res)
	if err != nil {
		return err
	}
	return ctx.Encode(res)
}

// Update decodes a resource from the Request, validates it and updates it in the database
func (h *BaseCRUD) Update(ctx *app.Context) error {
	res := h.NewResource(ctx)
	ID := ctx.URLParams.ByName(resourceIDParam)
	err := ctx.App.DB.Get(ID, res)
	if err != nil {
		return err
	}
	err = ctx.Decode(res)
	if err != nil {
		return err
	}
	err = h.checkAccess(ctx, res, "modify")
	if err != nil {
		return err
	}
	err = h.validator.Validate(database.ResourceActionUpdate, res)
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
func (h *BaseCRUD) Delete(ctx *app.Context) error {
	res := h.NewResource(ctx)
	ID := ctx.URLParams.ByName(resourceIDParam)
	err := ctx.App.DB.Get(ID, res)
	if err != nil {
		return err
	}
	err = h.checkAccess(ctx, res, "delete")
	if err != nil {
		return err
	}
	return ctx.App.DB.Delete(res)
}

// List retrieves a list of Resources from the database, and encodes it to the ResponseWriter
func (h *BaseCRUD) List(ctx *app.Context) error {
	list := h.NewResourceList(ctx)
	query, err := h.queryFromURL(ctx)
	if err != nil {
		return err
	}
	err = ctx.App.DB.FindMany(list, query)
	if err != nil {
		return err
	}
	err = expandResourceList(ctx, list)
	if err != nil {
		return err
	}
	return ctx.Encode(list)
}

func (h *BaseCRUD) queryFromURL(ctx *app.Context) (database.Query, error) {
	getInt := func(field string) (int, error) {
		val := ctx.Request.URL.Query().Get(field)
		if val == "" {
			return 0, nil
		}
		return strconv.Atoi(val)
	}

	ret := database.Query{}
	limit, err := getInt("limit")
	if err != nil {
		return ret, err
	}
	if limit == 0 {
		limit = h.resourceLimit
	}
	skip, err := getInt("skip")
	if err != nil {
		return ret, err
	}
	sort := ctx.Request.URL.Query()["sort_by"]

	var where interface{}
	if h.ownerField != "" {
		where = database.Dict{h.ownerField: ctx.User.GetID()}
	}
	return database.NewQuery(where, limit, skip, sort...), nil
}

func (h *BaseCRUD) checkAccess(ctx *app.Context, res database.Resource, action string) error {
	allowed, err := h.accessChecker.UserAllowed(ctx, ctx.User, res, action)
	if err != nil {
		return err
	}
	if allowed {
		return nil
	}
	ctx.WriteHeader(http.StatusForbidden)
	return ErrNotAllowed
}

func expandResource(ctx *app.Context, res database.Resource) error {
	rels := ctx.Request.URL.Query()["expand"]
	if len(rels) == 0 {
		return nil
	}
	return ctx.App.DB.FetchRelated(res, rels...)
}

func expandResourceList(ctx *app.Context, list interface{}) error {
	resList, err := database.AsResourceList(list)
	if err != nil {
		return err
	}
	for _, res := range resList {
		err := expandResource(ctx, res)
		if err != nil {
			return err
		}
	}
	return nil
}
