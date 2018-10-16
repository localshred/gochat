package main

import (
	"fmt"
	"net/http"
)

// HTTPServer represents an HTTP Server, duh
type HTTPServer struct {
	Config *httpServerConfig
	Server *Server
}

func newHTTPServer(server *Server, config *httpServerConfig) *HTTPServer {
	return &HTTPServer{
		Config: config,
		Server: server,
	}
}

func (httpServer *HTTPServer) start() {
	config := httpServer.Config
	context := &Context{}
	logger, logFile, err := createLogger("http - ", config.LogFile)
	if err != nil {
		printErrorAndExit(err, -1)
	}
	defer logFile.Close()
	context.Logger = logger

	address := fmt.Sprintf("%s:%v", httpServer.Config.Host, httpServer.Config.Port)
	handler := &HTTPHandler{
		Channels: &httpServer.Server.Channels,
		Context:  context,
	}
	server := &http.Server{
		Addr:    address,
		Handler: http.HandlerFunc(handler.ServeHTTP),
	}
	defer server.Close()

	listeningMessage := fmt.Sprintf("HTTP Server listening on %s\n", address)
	fmt.Print(listeningMessage)
	context.Logger.Debug(listeningMessage)

	server.ListenAndServe()
}
