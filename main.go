package main

import (
	"fmt"
	"os"
)

func main() {
	config := readConfig().Server
	server := newServer(config)
	server.start()
}

func printErrorAndExit(err error, exitCode int) {
	fmt.Fprintln(os.Stderr, "Unable to start server: ", err)
	os.Exit(exitCode)
}
