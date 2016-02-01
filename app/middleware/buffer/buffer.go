package buffer

import (
	"bytes"
	"io"
	"net/http"

	"bitbucket.org/syb-devs/goth/app"
)

type bufferedWriter struct {
	statusCode int
	buff       io.ReadWriter
	http.ResponseWriter
}

func newBufferedWriter(w http.ResponseWriter) *bufferedWriter {
	return &bufferedWriter{
		ResponseWriter: w,
		buff:           bytes.NewBuffer(nil),
	}
}

// Write writes the given data in the buffer
func (b *bufferedWriter) Write(data []byte) (int, error) {
	return b.buff.Write(data)
}

// WriteHeader sets the status code for the buffered writer
func (b *bufferedWriter) WriteHeader(statusCode int) {
	b.statusCode = statusCode
}

// Release writes the data and status code in the ResponseWriter
func (b *bufferedWriter) Release() (int64, error) {
	if b.statusCode != 0 {
		b.ResponseWriter.WriteHeader(b.statusCode)
	}
	return io.Copy(b.ResponseWriter, b.buff)
}

// New returns a middleware for buffered writing
func New() app.Middleware {
	return func(h app.Handler) app.Handler {
		return app.HandlerFunc(
			func(ctx *app.Context) error {
				buffw := newBufferedWriter(ctx.ResponseWriter)
				defer buffw.Release()

				ctx.ResponseWriter = buffw
				return h.Serve(ctx)
			})
	}
}
