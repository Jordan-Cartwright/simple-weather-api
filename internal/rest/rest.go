package rest

import (
	"api/pkg/render"
	"net/http"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// Response returns a struct for forming http response json
type Response struct {
	Message string `json:"message"`
}

// Respond writes a json response to a given http.ResponseWriter.
func Respond(w http.ResponseWriter, statusCode int, value interface{}) {
	if err := render.JSON(w, statusCode, value); err != nil {
		_, fn, line, _ := runtime.Caller(1)
		log.Errorf("Failed to write data to http ResponseWriter [%s:%d]: %s", fn, line, err)
	}
}

// RespondErr writes json response containing an internal server error status and error message
func RespondErr(w http.ResponseWriter, err error) {
	Respond(w, http.StatusInternalServerError, &Response{Message: err.Error()})
}
