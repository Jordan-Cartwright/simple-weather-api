package main

import (
	"api/config"
	owm "api/internal/openweathermap"
	"api/internal/rest"
	"api/internal/util"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

const (
	layoutISO = "2006-01-02"
)

var (
	configFile = flag.StringP("config", "c", "", "(optional) absolute path to the api configuration file")
	cfg        *config.Config
)

type windData struct {
	Speed   float64 `json:"speed"`
	Bearing float64 `json:"bearing"`
}

type dailyData struct {
	Date        string `json:"date"`
	WeatherType string `json:"type"`
	Description string `json:"description"`
	Temp        struct {
		Low  float64 `json:"low"`
		High float64 `json:"high"`
	} `json:"temperature"`
}

type weatherData struct {
	Date        string       `json:"date"`
	WeatherType string       `json:"type"`
	Description string       `json:"description"`
	Temp        float64      `json:"temperature"`
	Wind        *windData    `json:"wind"`
	PrecipProb  float64      `json:"precip_prob"`
	Daily       []*dailyData `json:"daily"`
}

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
	// parse the query params
	latitude := r.FormValue("latitude")
	longitude := r.FormValue("longitude")

	if latitude == "" && longitude == "" {
		rest.Respond(w, http.StatusBadRequest, &rest.Response{Message: "Missing query parametes `latitude` and `longitude`"})
		return
	} else if latitude == "" {
		rest.Respond(w, http.StatusBadRequest, &rest.Response{Message: "Missing the latitude value"})
		return
	} else if longitude == "" {
		rest.Respond(w, http.StatusBadRequest, &rest.Response{Message: "Missing the longitude value"})
		return
	}

	// convert param values to float64 and error check value
	var lat, lon float64
	if s, err := strconv.ParseFloat(latitude, 64); err == nil {
		lat = s
	} else {
		rest.Respond(w, http.StatusBadRequest, &rest.Response{Message: fmt.Sprintf("'%v' is an invalid latitude value", latitude)})
		return
	}
	if s, err := strconv.ParseFloat(longitude, 64); err == nil {
		lon = s
	} else {
		rest.Respond(w, http.StatusBadRequest, &rest.Response{Message: fmt.Sprintf("'%v' is an invalid longitude value", longitude)})
		return
	}

	// do the things to get weather data from chosen api
	weatherForecast, err := getWeatherForecast(lat, lon)
	if err != nil {
		rest.RespondErr(w, err)
		return
	}

	// respond with the data
	rest.Respond(w, http.StatusOK, weatherForecast)
}

// getWeatherForecast fetches the 7 day weather forecast with a given latitude and longitude
func getWeatherForecast(latitude float64, longitude float64) (*weatherData, error) {
	forecast, err := getWeatherData(latitude, longitude)
	if err != nil {
		return nil, err
	}

	// Format the raw data into the desired formatted response
	dailyForecastData := prepareDailyForecastData(forecast)

	data := &weatherData{
		Date:        formatUnixTime(forecast.Current.Dt, layoutISO),
		WeatherType: strings.Title(forecast.Current.Weather[len(forecast.Current.Weather)-1].Main),
		Description: strings.Title(forecast.Current.Weather[len(forecast.Current.Weather)-1].Description),
		Temp:        forecast.Current.Temperature,
		Wind: &windData{
			Speed:   forecast.Current.WindSpeed,
			Bearing: forecast.Current.WindDeg,
		},
		PrecipProb: forecast.Daily[0].Pop,
		Daily:      dailyForecastData,
	}

	return data, nil
}

// getWeatherData uses the openweathermap package to request the raw data from openweathermap.org
func getWeatherData(latitude float64, longitude float64) (*owm.OneCallData, error) {
	coordinates := &owm.Coordinates{
		Latitude:  latitude,
		Longitude: longitude,
	}
	onecall, err := owm.NewOneCall("F", "en", cfg.APIKey, []owm.ExcludeOption{owm.ExcludeMinutely, owm.ExcludeHourly, owm.ExcludeAlerts})
	if err != nil {
		return nil, err
	}
	onecall.PerformOneCall(coordinates)

	return onecall, nil
}

// formatUnixTime formats a given unix timestamp integer into
// a specified time string format
func formatUnixTime(unixTime int, format string) string {
	tm := time.Unix(int64(unixTime), 0)
	return tm.Format(format)
}

// prepareDailyForecastData creates the daily JSON for the weather endpoint response
func prepareDailyForecastData(forecast *owm.OneCallData) []*dailyData {
	var forecast7Day []*dailyData
	for _, day := range forecast.Daily {
		dailyInfo := &dailyData{
			Date:        formatUnixTime(day.Dt, layoutISO),
			WeatherType: strings.Title(day.Weather[len(day.Weather)-1].Main),
			Description: strings.Title(day.Weather[len(day.Weather)-1].Description),
			Temp: struct {
				Low  float64 "json:\"low\""
				High float64 "json:\"high\""
			}{
				Low:  day.Temp.Min,
				High: day.Temp.Max,
			},
		}
		forecast7Day = append(forecast7Day, dailyInfo)
	}
	return forecast7Day
}
