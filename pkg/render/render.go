package render

import (
	"encoding/json"
	"io"
	"net/http"
)

const contentTypeJSON = "application/json"

// JSON encodes a given value using the standard json package and writes
// the encoding output to a given writer. If the writer implements the
// http.ResponseWriter interface, then this function will also set the JSON
// content-type header. Status will be used when w is a http.ResponseWriter and
// must be a valid http status code.
func JSON(w io.Writer, status int, value interface{}) error {
	if hw, ok := w.(http.ResponseWriter); ok {
		hw.Header().Set("Content-Type", contentTypeJSON)
		hw.WriteHeader(status)
	}

	return json.NewEncoder(w).Encode(value)
}
