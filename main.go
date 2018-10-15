package main

import (
	"fmt"
	"os"
)

func main() {
	config := readConfig()
	telnetServer := &TelnetServer{
		Config: config.Server,
	}
	telnetServer.start()
}

func printErrorAndExit(err error, exitCode int) {
	fmt.Fprintln(os.Stderr, "Unable to start server: ", err)
	os.Exit(exitCode)
}
