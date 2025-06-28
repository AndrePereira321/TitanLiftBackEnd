package server

import (
	"github.com/gofiber/fiber/v3"
	"titan-lift/internal/config"
	"titan-lift/internal/database"
	"titan-lift/internal/logger"
	"titan-lift/internal/server_error"
)

type Server struct {
	fiber  *fiber.App
	db     *database.Database
	logger *logger.Logger
}

func New(serverConfig *config.ServerConfig) (*Server, error) {
	serverLogger, err := getServerLogger(serverConfig)
	if err != nil {
		return nil, server_error.Wrap("SERVER_INIT", "failed creating server logger", err)
	}

	db, err := database.New(serverConfig)
	if err != nil {
		return nil, server_error.Wrap("SERVER_INIT", "failed creating database", err)
	}

	return &Server{
		fiber:  nil,
		db:     db,
		logger: serverLogger,
	}, nil
}

func (s *Server) Close() error {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			return server_error.Wrap("SERVER_CLOSE", "failed closing database", err)
		}
	}
	return nil
}

func getServerLogger(serverConfig *config.ServerConfig) (*logger.Logger, error) {
	level := serverConfig.Logging().ServerLogLevel()
	dir := serverConfig.Logging().LogDir()
	return logger.New("SERVER", level, dir)
}
