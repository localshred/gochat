package main

var configFile = "./configs/config.json"

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
	// if rawJSON, err := ioutil.ReadFile(configFile); err != nil {
	// panic(err)
	// }
	return defaultConfig
}
