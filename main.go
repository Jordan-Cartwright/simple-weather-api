package main

import (
	"api/config"
	"api/internal/rest"
	"api/internal/util"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	configFile = flag.StringP("config", "c", "", "(optional) absolute path to the api configuration file")
	cfg        *config.Config
)

func main() {
	flag.Parse()

	cfg = config.NewConfig(*configFile)

	util.InitializeLogging(os.Stderr, cfg.Log.Level, cfg.Log.Format)

	port := cfg.Port
	addr := fmt.Sprintf(":%v", port)

	log.Infof("APP is listening on port: %s", port)
	log.Fatal(http.ListenAndServe(addr, handler()))
}

func handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/api/v1/weather", GetForecast)
	r.HandleFunc("/api/v1/ping", GetStatus)

	return r
}

// GetStatus returns a JSON response when the server is running
func GetStatus(w http.ResponseWriter, r *http.Request) {
	message := &rest.Response{
		Message: "pong",
	}
	rest.Respond(w, http.StatusOK, message)
}

// GetForecast returns a JSON response with a 7 day weather forecast
func GetForecast(w http.ResponseWriter, r *http.Request) {
	rest.Respond(w, http.StatusNotImplemented, rest.Response{Message: "endpoint not implemented"})
}
