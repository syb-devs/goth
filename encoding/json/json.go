package json

import (
	"encoding/json"
	"io"
	"net/http"
)

// Codec encodes and decodes data from/to JSON text
type Codec struct{}

// Encode marshals a resource (r) as JSON data in the given writer (w)
func (c Codec) Encode(w io.Writer, res interface{}) error {
	if rw, ok := w.(http.ResponseWriter); ok {
		if h := rw.Header().Get("Content-Type"); h == "" {
			rw.Header().Set("Content-Type", "application/json")
		}
	}
	return json.NewEncoder(w).Encode(res)
}

// Decode unmarshals JSON data read from the Reader (r) into res
func (c Codec) Decode(r io.Reader, res interface{}) error {
	return json.NewDecoder(r).Decode(res)
}
