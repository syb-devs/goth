package encoding

import (
	"io"
)

// Codec interface is used to encode and decode data between the app
// and the transport layers (HTTP, websockets)
type Codec interface {
	Encode(io.Writer, interface{}) error
	Decode(io.Reader, interface{}) error
}
