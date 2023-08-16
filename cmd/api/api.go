package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/cyclimse/fediverse-blahaj/internal/api/controller"
	api "github.com/cyclimse/fediverse-blahaj/internal/api/v1"
	"github.com/cyclimse/fediverse-blahaj/internal/business"
)

func main() {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, "postgres://fediverse:fediverse@localhost:5432/fediverse")
	if err != nil {
		panic(err)
	}
	defer dbpool.Close()

	b := business.New(dbpool)

	e := echo.New()

	// Add some middleware
	// e.Use(middleware.Recover())
	e.Use(middleware.Timeout())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TODO: Change this to the actual domain
		// This is only for development
		AllowOrigins: []string{"http://localhost:5173"},
	}))

	api.RegisterHandlersWithBaseURL(e, controller.NewAPIController(b), "/api/v1")

	e.Logger.Fatal(e.Start(":8080"))
}
