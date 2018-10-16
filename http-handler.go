package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

// HTTPHandler a type for responding to HTTP requests
type HTTPHandler struct {
	Channels *map[string]*Channel
	Context  *Context
}

type handlerFunc func(*HTTPHandler, http.ResponseWriter, *http.Request) (n int, statusCode int)

var endpoints = map[string]map[string]handlerFunc{
	"GET": map[string]handlerFunc{
		"/channels(\\?.*|$)":              handleGetChannels,
	},
}

// ChannelsJSON represents a list of channel names to be responded as JSON
type ChannelsJSON struct {
	Channels []string `json:"channels"`
}

func handleGetChannels(handler *HTTPHandler, response http.ResponseWriter, request *http.Request) (n int, statusCode int) {
	channels := listChannels(*handler.Channels)
	payload := map[string][]string{
		"channels": channels,
	}

	contentType := "application/json"
	statusCode = 200
	bytes, err := json.Marshal(payload)
	if nil != err {
		contentType = "text/plain"
		statusCode = 404
		bytes = []byte(err.Error())
	}
	return writeResponse(response, statusCode, contentType, bytes)
}

func handleNotFound(handler *HTTPHandler, response http.ResponseWriter, request *http.Request) (n int, statusCode int) {
	return writeResponse(response, 404, "text/plain", []byte("Not Found"))
}

func (handler *HTTPHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	uri := request.RequestURI
	method := request.Method
	methodRouter, ok := endpoints[method]
	if ok {
		for pattern, routeHandler := range methodRouter {
			if matched, _ := regexp.MatchString(pattern, uri); matched {
				n, statusCode := routeHandler(handler, response, request)
				handler.Context.Logger.Debugf("%v %v %s %s", statusCode, n, method, uri)
				return
			}
		}
	}

	n, statusCode := handleNotFound(handler, response, request)
	handler.Context.Logger.Warnf("%v %v %s %s", statusCode, n, method, uri)
}

func writeResponse(response http.ResponseWriter, statusCode int, contentType string, bytes []byte) (int, int) {
	n := len(bytes)
	response.Header().Set("Content-Type", contentType)
	response.WriteHeader(statusCode)
	response.Write(bytes)
	return n, statusCode
}
