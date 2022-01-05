package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/syb-devs/goth/database"

	"golang.org/x/net/websocket"
)

var ErrNoAuth = errors.New("authentication needed")

// Compile-time interface check
var _ Handler = (*WSServer)(nil)

// WSHandlerFunc handles a Websocket WSEvent in a given WSConnection
type WSHandlerFunc func(*WSConn, *WSEvent) error

type bindFunc func(*WSServer) error

var bindFns = make([]bindFunc, 0)

func RegisterWSBindFunc(f bindFunc) {
	bindFns = append(bindFns, f)
}

// WSEvent represents a single unit of information exchanged between the websocket client and WSServer
type WSEvent struct {
	Name string          `json:"event"`
	Data json.RawMessage `json:"data"`
}

// DecodeData decodes the event raw data in the given struct
func (e *WSEvent) DecodeData(dest interface{}) error {
	err := json.Unmarshal(e.Data, &dest)
	if err != nil {
		return err
	}
	return nil
}

// WSEvent represents a single unit of information exchanged between the websocket client and WSServer
type WSEventOut struct {
	Name string      `json:"event"`
	Data interface{} `json:"data"`
}

// WSConn represents a websocket WSConnection
type WSConn struct {
	ws     *websocket.Conn
	srv    *WSServer
	DBConn database.Connection
	IsAuth bool
	UserID string
}

// Dispatch fires the handler for a given WSEvent
func (c *WSConn) Dispatch(e *WSEvent) error {
	err := c.dispatch(e)
	if err != nil {
		return err
	}
	if e.Name == "auth" {
		return c.dispatch(&WSEvent{Name: "init"})
	}
	return nil
}

func (c *WSConn) dispatch(WSEvent *WSEvent) error {
	hs := c.srv.ehs[WSEvent.Name]
	if len(hs) == 0 {
		return nil
	}

	for _, h := range hs {
		err := h(c, WSEvent)
		if err != nil {
			return err
		}
	}
	return nil
}

// Listen locks on the websocket WSConnection and reads incoming WSEvents
func (c *WSConn) Listen() error {
	//TODO add exit conditions
	for {
		var e WSEvent
		err := websocket.JSON.Receive(c.ws, &e)
		if err != nil {
			log.Printf("websocket: error receiving: %v\n", err)
			return err
		}
		if !c.IsAuth && e.Name != "auth" {
			return ErrNoAuth
		}

		err = c.Dispatch(&e)
		if err != nil {
			return err
		}
	}
}

// SendEvent sends the WSEvent along the WSConnection
func (c *WSConn) SendEvent(e *WSEventOut) error {
	websocket.JSON.Send(c.ws, e)
	return nil
}

// Close performs cleanup tasks on connection end and triggers binded handlers
func (c *WSConn) Close() error {
	c.Dispatch(&WSEvent{Name: "close"})
	c.DBConn.Close()
	return c.ws.Close()
}

// WSServer handles websocket WSConnections
type WSServer struct {
	ehs map[string][]WSHandlerFunc
	app *App
}

// NewWSServer creates a new websocket WSServer
func NewWSServer(app *App) *WSServer {
	return &WSServer{
		ehs: make(map[string][]WSHandlerFunc, 0),
		app: app,
	}
}

// LoadEventHandlers calls the registered event-binding functions
func (s *WSServer) LoadEventHandlers() {
	for _, f := range bindFns {
		if err := f(s); err != nil {
			panic(err)
		}
	}
}

// Bind registers an WSEvent name to the corresponding handler
func (s *WSServer) Bind(WSEvent string, h WSHandlerFunc) {
	s.ehs[WSEvent] = append(s.ehs[WSEvent], h)
}

// ServeHTTP implements the Handler interface
func (s *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request, _ *Context) error {
	websocket.Handler(s.handleWSConn).ServeHTTP(w, r)
	return nil
}

// Echo the data received on the WebSocket.
func (s *WSServer) handleWSConn(ws *websocket.Conn) {
	log.Println("WS connection started")
	conn := &WSConn{
		ws:     ws,
		srv:    s,
		DBConn: s.app.DB.Copy(),
	}
	defer conn.Close()

	conn.Dispatch(&WSEvent{Name: "open"})
	if err := conn.Listen(); err != nil {
		log.Printf("websocket error: %v", err)
	}
}
