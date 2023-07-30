package main

import (
	"context"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"

	"github.com/cyclimse/fediverse-blahaj/internal/business"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/cyclimse/fediverse-blahaj/internal/orchestrator"
)

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://fediverse:fediverse@localhost:5432/fediverse")
	if err != nil {
		panic(err)
	}

	b := business.New(conn)
	o := orchestrator.New()

	// search for peers for 10 seconds
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// create a channel to receive the results
	results := make(chan models.FediverseServer, 100)

	// run the crawler and the business logic in parallel
	// using errgroup to exit early if one of them fails
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(2)

	g.Go(func() error {
		return o.Crawl(ctx, results)
	})
	g.Go(func() error {
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
