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
		{name: "ping", path: "/api/v1/weather", status: http.StatusNotImplemented},
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

func TestGetForcast(t *testing.T) {}
