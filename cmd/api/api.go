package main

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/cyclimse/fediverse-blahaj/internal/api/controller"
	api "github.com/cyclimse/fediverse-blahaj/internal/api/v1"
	"github.com/cyclimse/fediverse-blahaj/internal/business"
)

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://fediverse:fediverse@localhost:5432/fediverse")
	if err != nil {
		panic(err)
	}

	b := business.New(conn)

	e := echo.New()

	// Add some middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Timeout())

	api.RegisterHandlersWithBaseURL(e, controller.NewAPIController(b), "/api/v1")

	e.Logger.Fatal(e.Start(":8080"))
}
