package main

import (
	"api/internal/rest"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouting(t *testing.T) {
	testCases := []struct {
		name   string
		path   string
		status int
	}{
		{name: "ping", path: "/api/v1/ping", status: http.StatusOK},
		{name: "weather", path: "/api/v1/weather", status: http.StatusBadRequest},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(handler())
			defer srv.Close()

			response, err := http.Get(fmt.Sprintf("%s%s", srv.URL, tc.path))
			if err != nil {
				t.Fatalf("could not send GET request: %v", err)
			}

			if response.StatusCode != tc.status {
				t.Errorf("expected status: %v; actual status: %v", tc.status, response.StatusCode)
			}
		})
	}
}

func TestGetStatus(t *testing.T) {
	testCases := []struct {
		name     string
		status   int
		response rest.Response
	}{
		{
			name: "ping", status: http.StatusOK, response: rest.Response{Message: "pong"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "localhost:8080/api/v1/ping", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			recorder := httptest.NewRecorder()

			GetStatus(recorder, request)

			response := recorder.Result()
			defer response.Body.Close()

			// Check the status code
			if response.StatusCode != tc.status {
				t.Fatalf("expected status: %v, actual status: %v", tc.status, response.StatusCode)
			}

			// Check the body of the response
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			actualResponse := bytes.TrimSpace(body)
			expectedResponse, err := json.Marshal(tc.response)
			if err != nil {
				t.Fatalf("could not unmarshal expected response: %v", err)
			}

			if string(actualResponse) != string(expectedResponse) {
				t.Fatalf("expected response: %v; actual response: %v", expectedResponse, actualResponse)
			}

		})
	}
}

func TestGetForcast(t *testing.T) {
	testCases := []struct {
		name     string
		lat      string
		long     string
		status   int
		response rest.Response
		err      string
	}{
		// {name: "latitude and longitude", lat: "33.7984", long: "-84.3883", status: http.StatusOK, response: rest.Response{Message: "example data"}},
		{name: "no latitude or longitude", lat: "", long: "", status: http.StatusBadRequest, response: rest.Response{Message: "Missing query parametes `latitude` and `longitude`"}},
		{name: "only latitude", lat: "33.7984", long: "", status: http.StatusBadRequest, response: rest.Response{Message: "Missing the longitude value"}},
		{name: "only longitude", lat: "", long: "-84.3883", status: http.StatusBadRequest, response: rest.Response{Message: "Missing the latitude value"}},
		{name: "invalid latitude", lat: "thing", long: "-84.3883", status: http.StatusBadRequest, response: rest.Response{Message: "'thing' is an invalid latitude value"}},
		{name: "invalid longitude", lat: "33.7984", long: "thing", status: http.StatusBadRequest, response: rest.Response{Message: "'thing' is an invalid longitude value"}},
		{name: "invalid latitude and longitude", lat: "thing1", long: "thing2", status: http.StatusBadRequest, response: rest.Response{Message: "'thing1' is an invalid latitude value"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "localhost:8080/weather?latitude="+tc.lat+"&longitude="+tc.long, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			recorder := httptest.NewRecorder()

			GetForecast(recorder, req)

			response := recorder.Result()
			defer response.Body.Close()

			// Check the status code
			if response.StatusCode != tc.status {
				t.Fatalf("expected status: %v, actual status: %v", tc.status, response.StatusCode)
			}

			// Check the body of the response
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			actualResponse := bytes.TrimSpace(body)
			expectedResponse, err := json.Marshal(tc.response)
			if err != nil {
				t.Fatalf("could not unmarshal expected response: %v", err)
			}

			if string(actualResponse) != string(expectedResponse) {
				t.Fatalf("expected response: %v; actual response: %v", string(expectedResponse), string(actualResponse))
			}

		})
	}
}

func TestFormatUnixTime(t *testing.T) {
	testCases := []struct {
		name          string
		unixTime      int
		format        string
		expectedValue string
	}{
		{name: "2021", unixTime: 1624683380, format: "2006-01-02", expectedValue: "2021-06-26"},
		{name: "2012", unixTime: 1356170165, format: "2006-01-02", expectedValue: "2012-12-22"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value := formatUnixTime(tc.unixTime, tc.format)

			if value != tc.expectedValue {
				t.Errorf("expected time: %v; actual time: %v", tc.expectedValue, value)
			}
		})
	}
}
