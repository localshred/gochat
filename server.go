package main

func main() {
	config := readConfig()
type chatServerConfig struct {
	Host    string
	Port    int
	LogFile string
}

func readConfig() *chatServerConfig {
	defaultConfig := &chatServerConfig{
		Host:    "localHost",
		LogFile: "logs/server.log",
		Port:    5555,
	}
	return defaultConfig
}
