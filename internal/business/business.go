package business

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/cyclimse/fediverse-blahaj/internal/db"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(conn *pgxpool.Pool) *Business {
	return &Business{
		conn:    conn,
		queries: db.New(conn),
	}
}

type Business struct {
	conn    *pgxpool.Pool
	queries *db.Queries
}

// Run runs the business logic.
// Handles the results from the crawler, but is not aware of the crawler logic.
func (b *Business) Run(ctx context.Context, servers chan models.FediverseServer) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case server, ok := <-servers:
			if !ok {
				return nil
			}
			slog.InfoCtx(ctx, "received server", "server", server)
			// check if the server is already in the db
			var s db.Server
			s, err := b.queries.GetSeverByDomain(ctx, server.Domain)
			if err != nil {
				if err != pgx.ErrNoRows {
					return err
				}
				// server is not in the db, add it
				s, err = b.AddServer(ctx, server)
				if err != nil {
					slog.ErrorCtx(ctx, "failed to add server", "error", err)
					return err
				}
			}
			// add the crawl to the server
			if server.CrawlErr != nil {
				err = b.AddFailedCrawlToServer(ctx, server, s)
			} else {
				err = b.AddCompletedCrawlToServer(ctx, server, s)
			}
			if err != nil {
				slog.ErrorCtx(ctx, "failed to add crawl to server", "error", err)
				return err
			}
		}
	}
}
