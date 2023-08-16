package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"

	"github.com/cyclimse/fediverse-blahaj/internal/business"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/cyclimse/fediverse-blahaj/internal/orchestrator"
)

func main() {
	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, "postgres://fediverse:fediverse@localhost:5432/fediverse")
	if err != nil {
		panic(err)
	}
	defer dbpool.Close()

	b := business.New(dbpool)
	o := orchestrator.New(business.BlockedDomains)

	// create a channel to receive the results
	results := make(chan models.FediverseServer, 100)

	// run the crawler and the business logic in parallel
	// using errgroup to exit early if one of them fails
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(2)

	g.Go(func() error {
		// this will run until the context is cancelled
		defer close(results)
		// create a context with a timeout for the crawl
		crawlCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
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
			return
		}
		panic(err)
	}
}
