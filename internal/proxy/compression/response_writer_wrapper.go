package compression

import (
	"io"
	"net/http"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	io.Writer
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	return rw.Writer.Write(b)
}
