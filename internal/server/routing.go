package server

import "titan-lift/internal/server/routes"

type RouteHandler func(ctx *routes.RouteContext) error

func (s *Server) RegisterRoutes() {
	s.Post("/register", routes.RegisterUser)
}
