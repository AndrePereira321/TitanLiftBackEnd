package routes

import (
	"github.com/gofiber/fiber/v3"
	"titan-lift/internal/database"
)

type UserContext struct {
	userId    uint64
	sessionId uint64
}
type RouteContext struct {
	ctx         *fiber.Ctx
	userContext *UserContext
}

func GetRouteContext(ctx *fiber.Ctx, db *database.Database) (*RouteContext, error) {
	//TODO Retrieve user context from DB
	return &RouteContext{
		ctx:         ctx,
		userContext: nil,
	}, nil
}
