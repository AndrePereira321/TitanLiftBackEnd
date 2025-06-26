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
	logger, err := getServerLogger(serverConfig)
	if err != nil {
		panic(err)
	}

	defer logger.Close()

	fmt.Println(logger)
	fmt.Println(serverConfig)
}

func getServerLogger(config *config.ServerConfig) (*Logger, error) {
	level := config.Logging().ServerLogLevel()
	dir := config.Logging().LogDir()
	return NewLogger("SERVER", level, dir)
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
