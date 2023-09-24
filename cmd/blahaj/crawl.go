package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"golang.org/x/sync/errgroup"

	"github.com/cyclimse/fediverse-blahaj/internal/business"
	"github.com/cyclimse/fediverse-blahaj/internal/config"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/cyclimse/fediverse-blahaj/internal/orchestrator"
)

const (
	MaximumConcurrentCrawls = 100
)

type CrawlCmd struct {
	// To be set to the Container Timeout in production
	Duration     time.Duration `help:"Duration of the crawl." default:"5m" env:"CRAWL_DURATION"`
	CrawlerCount int           `help:"Number of crawlers." default:"2" env:"CRAWLER_COUNT"`

	EntryPointServerPort int `help:"Port to listen on for the entry point server." default:"8081" env:"PORT"`
}

func (cmd *CrawlCmd) Run(cmdContext *Context) error {
	cfg := cmdContext.Config

	if cfg.Environment == config.EnvironmentDevelopment {
		return cmd.RunCrawl(cmdContext)
	}

	// in production, we run the entry point server
	return cmd.RunEntryPointServer(cmdContext)
}

func (cmd *CrawlCmd) RunCrawl(cmdContext *Context) error {
	cfg := cmdContext.Config
	cfg.SetDevelopmentDefaults()

	dbpool, err := pgxpool.New(cmdContext.Ctx, cfg.PgConn)
	if err != nil {
		return err
	}
	defer dbpool.Close()

	b := business.New(dbpool)

	seeds, err := b.GetCrawlerSeedDomains(cmdContext.Ctx, 50)
	if err != nil {
		return err
	}

	o := orchestrator.New(orchestrator.OrchestratorConfig{
		NumCrawlers:      cmd.CrawlerCount,
		BlockedDomains:   business.BlockedDomains,
		SeedDomains:      seeds,
		CrawlTimeout:     cmd.Duration,
		CrawlerUserAgent: fmt.Sprintf("blahaj/%s", cmdContext.Version),
	})

	// create a channel to receive the results
	results := make(chan models.Crawl, MaximumConcurrentCrawls)

	// run the crawler and the business logic in parallel
	// using errgroup to exit early if one of them fails
	g, ctx := errgroup.WithContext(cmdContext.Ctx)
	g.SetLimit(2)

	g.Go(func() error {
		// this will run until the context is cancelled
		defer close(results)
		// create a context with a timeout for the crawl
		crawlCtx, cancel := context.WithTimeout(ctx, cmd.Duration)
		defer cancel()

		slog.InfoContext(crawlCtx, "starting crawl", "seeds", len(seeds), "crawlers", cmd.CrawlerCount)

		err := o.Crawl(crawlCtx, results)
		if err != nil && crawlCtx.Err() != nil {
			// we do not want to return an error because
			// it would cancel the other goroutine in the group
			return nil
		}
		return err
	})
	g.Go(func() error {
		// this will run until the results channel is closed
		return b.Run(ctx, results)
	})

	// wait for the goroutines to finish
	if err := g.Wait(); err != nil {
		if ctx.Err() != nil {
			// context error, this is expected
			// the context is cancelled when the timeout is reached
			slog.InfoContext(ctx, "context error", "err", err)
			return nil
		}
		slog.ErrorContext(ctx, "error", "err", err)
		return err
	}

	return nil
}

// This server is used as an entry point to start the crawl.
// It used in production because of the way Scaleway Serverless Containers work.
func (cmd *CrawlCmd) RunEntryPointServer(cmdContext *Context) error {
	e := echo.New()
	e.HideBanner = true

	atomicIsAlreadyCrawling := uint32(0)
	crawlingCtx, cancel := context.WithCancel(cmdContext.Ctx)
	defer cancel()

	e.GET("/", func(c echo.Context) error {
		switch atomic.LoadUint32(&atomicIsAlreadyCrawling) {
		case 0:
			return c.String(200, "Not crawling.")
		default:
			return c.String(200, "Crawling.")
		}
	})

	e.POST("/stop", func(c echo.Context) error {
		// check if we are already crawling
		if atomic.LoadUint32(&atomicIsAlreadyCrawling) == 0 {
			return c.String(200, "Not crawling, nothing to do.")
		}

		// cancel the context to stop the crawl
		cancel()

		return c.String(200, "Stopped crawling.")
	})

	e.POST("/", func(c echo.Context) error {
		// check if we are already crawling
		if atomic.LoadUint32(&atomicIsAlreadyCrawling) >= 1 {
			return c.String(200, "Already crawling, nothing to do.")
		}

		// set the flag to prevent multiple crawls
		atomic.StoreUint32(&atomicIsAlreadyCrawling, 1)

		// run the crawl
		go func() {
			defer func() {
				// recover from panic
				if r := recover(); r != nil {
					slog.ErrorContext(cmdContext.Ctx, "panic", "panic", r)
				}
				// reset the flag when the crawl is finished
				atomic.StoreUint32(&atomicIsAlreadyCrawling, 0)
			}()

			// Swap the context to allow cancelling the crawl
			cmdContext.Ctx = crawlingCtx
			err := cmd.RunCrawl(cmdContext)
			if err != nil {
				// TODO: add a Slack notification? Or some sort of AlertManager integration?
				slog.ErrorContext(cmdContext.Ctx, "error", "err", err)
			}
		}()

		return c.String(200, "Started crawling.")
	})

	return e.Start(fmt.Sprintf(":%d", cmd.EntryPointServerPort))
}
