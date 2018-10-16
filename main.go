package main

import (
	"fmt"
	"os"
)

func main() {
	config := readConfig()
	server := newServer(config.Server)

	httpServer := newHTTPServer(server, config.HTTPServer)
	go httpServer.start()

	server.start()
}

func printErrorAndExit(err error, exitCode int) {
	fmt.Fprintln(os.Stderr, "Unable to start server: ", err)
	os.Exit(exitCode)
}
