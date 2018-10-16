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

func writeResponse(response http.ResponseWriter, statusCode int, contentType string, bytes []byte) (int, int) {
	n := len(bytes)
	response.Header().Set("Content-Type", contentType)
	response.Write(bytes)
	response.WriteHeader(statusCode)
	return n, statusCode
}

func routeNotFound(handler *HTTPHandler, response http.ResponseWriter, request *http.Request) (n int, statusCode int) {
	return writeResponse(response, 404, "text/plain", []byte("Not Found"))
}

var endpoints = map[string]map[string]handlerFunc{
	"GET": map[string]handlerFunc{
	},
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

	n, statusCode := routeNotFound(handler, response, request)
	handler.Context.Logger.Warnf("%v %v %s %s", statusCode, n, method, uri)
}
