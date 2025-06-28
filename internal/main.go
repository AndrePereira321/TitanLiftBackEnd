package main

import (
	"os"
	"titan-lift/internal/config"
	"titan-lift/internal/database"
	"titan-lift/internal/logger"
	"titan-lift/internal/server_error"
)

func main() {
	serverConfig, err := getConfig()
	if err != nil {
		panic(err)
	}

	appLogger, err := logger.New("APP", "TRACE", serverConfig.Logging().LogDir())
	if err != nil {
		panic(err)
	}
	defer func(appLogger *logger.Logger) {
		_ = appLogger.Close()
	}(appLogger)

	db, err := database.New(serverConfig)
	if err != nil {
		appLogger.FatalEvent().Err(err).Msg("failed to initialize database")
	}
	defer func(database *database.Database) {
		_ = database.Close()
	}(db)
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
