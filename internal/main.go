package main

import (
	"fmt"
	"os"
	"titan-lift/internal/config"
	"titan-lift/internal/server_error"
)

func main() {
	serverConfig, err := getConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println(serverConfig)
}

func getConfig() (*config.ServerConfig, error) {
	if len(os.Args) < 2 {
		return nil, server_error.New("INIT", "missing configuration file path as 1st argument")
	}
	file, err := os.ReadFile(os.Args[1])
	if err != nil {
		return nil, server_error.Wrap("INIT", "error when reading config file", err)
	}
	return config.GetServerConfig(file)
}
