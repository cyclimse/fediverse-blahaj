package main

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/cyclimse/fediverse-blahaj/internal/api/controller"
	api "github.com/cyclimse/fediverse-blahaj/internal/api/v1"
	"github.com/cyclimse/fediverse-blahaj/internal/business"
	"github.com/cyclimse/fediverse-blahaj/internal/config"
)

type APICmd struct {
	Port int `help:"Port to listen on." default:"8080" env:"PORT"`
}

func (cmd *APICmd) Run(cmdContext *Context) error {
	cfg := cmdContext.Config
	cfg.SetDevelopmentDefaults()

	dbpool, err := pgxpool.New(cmdContext.Ctx, cfg.PgConn)
	if err != nil {
		return err
	}
	defer dbpool.Close()

	b := business.New(dbpool)

	e := echo.New()

	if cmdContext.Debug {
		e.Debug = true
	}

	// Add some middlewares
	if cfg.Environment != config.EnvironmentDevelopment {
		e.Use(middleware.Recover())
		e.Use(middleware.Gzip())
		// On production, we run behind a reverse proxy
		e.IPExtractor = echo.ExtractIPFromXFFHeader()
		// Add security headers
		e.Use(middleware.Secure())
		// Also force HTTPS
		e.Pre(middleware.HTTPSRedirect())
		// Finally hide the banner to make it easier to parse the logs
		e.HideBanner = true
	}

	e.Use(middleware.Timeout())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// For practical reasons, we split the API and the frontend.
		// This means we need to allow CORS for the frontend.
		AllowOrigins: []string{cfg.FrontendURL},
	}))

	api.RegisterHandlersWithBaseURL(e, controller.NewAPIController(b), "/api/v1")

	return e.Start(fmt.Sprintf(":%d", cmd.Port))
}
