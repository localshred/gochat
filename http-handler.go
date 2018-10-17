package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// HTTPHandler a type for responding to HTTP requests
type HTTPHandler struct {
	Channels   *map[string]*Channel
	Context    *Context
	Dispatcher *chan *Message
}

type handlerFunc func(*HTTPHandler, map[string]string, http.ResponseWriter, *http.Request) (n int, statusCode int)

var endpoints = map[string]map[string]handlerFunc{
	"GET": map[string]handlerFunc{
		"/channels/:channelName/messages": handleGetChannelMessages,
		"/channels":                       handleGetChannels,
	},
	"POST": map[string]handlerFunc{
		"/channels/:channelName/messages": handlePostChannelMessage,
	},
}

// ChannelsJSON represents a list of channel names to be responded as JSON
type ChannelsJSON struct {
	Channels []string `json:"channels"`
}

// ChannelMessagesJSON represents a list of channel messages to be responded as JSON
type ChannelMessagesJSON struct {
	Messages []*Message `json:"messages"`
}

// PostMessageRequestJSON represents the POST body as a JSON document when creating a message on a channel
type PostMessageRequestJSON struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// PostMessageResponseJSON represents the JSON response when creating a message on a channel
type PostMessageResponseJSON struct {
	Message *Message `json:"message"`
}

func getURLParams(template, uri string) (params map[string]string) {
	params = map[string]string{}
	if matched, _ := regexp.MatchString(".*/:[^/]+/?", template); !matched {
		return
	}
	templateFields := strings.Split(template, "/")
	uriFields := strings.Split(uri, "/")
	if len(templateFields) != len(uriFields) {
		return
	}

	for i := 0; i < len(templateFields); i++ {
		templateField := templateFields[i]
		if matched, _ := regexp.MatchString("^:[^/]+$", templateField); matched {
			name := templateField[1:]
			value := uriFields[i]
			params[name] = value
		}
	}
	return
}

func handleGetChannelMessages(handler *HTTPHandler, urlParams map[string]string, response http.ResponseWriter, request *http.Request) (n int, statusCode int) {
	channelName, ok := urlParams["channelName"]
	if !ok {
		return writeResponse(response, 404, "text/plain", []byte("Not Found"))
	}

	channel, ok := (*handler.Channels)[channelName]
	if !ok {
		return writeResponse(response, 404, "text/plain", []byte("Not Found"))
	}
	payload := &ChannelMessagesJSON{
		Messages: channel.Messages,
	}

	return writeJSONResponse(response, 200, payload)
}

func handleGetChannels(handler *HTTPHandler, urlParams map[string]string, response http.ResponseWriter, request *http.Request) (n int, statusCode int) {
	channels := listChannels(*handler.Channels)
	payload := &ChannelsJSON{
		Channels: channels,
	}

	return writeJSONResponse(response, 200, payload)
}

func handlePostChannelMessage(handler *HTTPHandler, urlParams map[string]string, response http.ResponseWriter, request *http.Request) (n int, statusCode int) {
	channelName, ok := urlParams["channelName"]
	if !ok {
		channelName = "general"
	}
	data := &PostMessageRequestJSON{}
	parseJSONBody(request.Body, data)

	channel := (*handler.Channels)[channelName]
	message := &Message{
		Channel: channel,
		Message: data.Message,
		Time:    time.Now().UTC(),
		User:    &User{Username: data.Username},
	}
	*handler.Dispatcher <- message

	payload := &PostMessageResponseJSON{Message: message}
	return writeJSONResponse(response, 200, payload)
}

func handleNotFound(handler *HTTPHandler, urlParams map[string]string, response http.ResponseWriter, request *http.Request) (n int, statusCode int) {
	return writeResponse(response, 404, "text/plain", []byte("Not Found"))
}

func matchesEndpointURI(endpoint, uri string) (matched bool) {
	matched = false
	re, err := regexp.Compile("/:[^/]+")
	if nil != err {
		return
	}
	patternString := re.ReplaceAllLiteralString(regexp.QuoteMeta(endpoint), "/[^/]+")
	re, err = regexp.Compile(fmt.Sprintf("^%s/?$", patternString))
	if nil != err {
		return
	}
	matched = re.MatchString(uri)
	return
}

func parseJSONBody(body io.Reader, into interface{}) error {
	scanner := bufio.NewScanner(body)
	scanner.Split(bufio.ScanLines)
	data := []byte{}
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}
	if err := json.Unmarshal(data, into); nil != err {
		return err
	}
	return nil
}

func (handler *HTTPHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	uri := request.RequestURI
	method := request.Method
	methodRouter, ok := endpoints[method]
	if ok {
		for routeTemplate, routeHandler := range methodRouter {
			if matched := matchesEndpointURI(routeTemplate, uri); matched {
				urlParams := getURLParams(routeTemplate, uri)
				n, statusCode := routeHandler(handler, urlParams, response, request)
				handler.Context.Logger.Debugf("%v %v %s %s", statusCode, n, method, uri)
				return
			}
		}
	}

	n, statusCode := handleNotFound(handler, map[string]string{}, response, request)
	handler.Context.Logger.Warnf("%v %v %s %s", statusCode, n, method, uri)
}

func writeJSONResponse(response http.ResponseWriter, statusCode int, payload interface{}) (int, int) {
	contentType := "application/json"
	bytes, err := json.Marshal(payload)
	if nil != err {
		contentType = "text/plain"
		statusCode = 404
		bytes = []byte(err.Error())
	}
	return writeResponse(response, statusCode, contentType, bytes)
}

func writeResponse(response http.ResponseWriter, statusCode int, contentType string, bytes []byte) (int, int) {
	n := len(bytes)
	response.Header().Set("Content-Type", contentType)
	response.WriteHeader(statusCode)
	response.Write(bytes)
	return n, statusCode
}
