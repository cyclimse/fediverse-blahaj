package business

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/cyclimse/fediverse-blahaj/internal/db"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/cyclimse/fediverse-blahaj/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func New(conn *pgx.Conn) *Business {
	return &Business{
		conn:    conn,
		queries: db.New(conn),
	}
}

type Business struct {
	conn    *pgx.Conn
	queries *db.Queries
}

// Run runs the business logic.
// Handles the results from the crawler, but is not aware of the crawler logic.
func (b *Business) Run(ctx context.Context, servers chan models.FediverseServer) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case server := <-servers:
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
			err = b.AddCrawlToServer(ctx, server, s)
			if err != nil {
				slog.ErrorCtx(ctx, "failed to add crawl to server", "error", err)
				return err
			}
		}
	}
}

func (b *Business) AddServer(ctx context.Context, server models.FediverseServer) (db.Server, error) {
	s, error := b.queries.CreateServer(ctx, db.CreateServerParams{
		Domain:       server.Domain,
		SoftwareName: pgtype.Text{String: server.SoftwareName, Valid: true},
	})
	if error != nil {
		return db.Server{}, error
	}
	return s, nil
}

func (b *Business) AddCrawlToServer(ctx context.Context, res models.FediverseServer, s db.Server) error {
	c, err := b.queries.CreateCrawl(ctx, db.CreateCrawlParams{
		ServerID:          s.ID,
		NumberOfPeers:     int32(len(res.Peers)),
		OpenRegistrations: res.OpenRegistrations,
		TotalUsers:        pgtype.Int4{Int32: int32(utils.IntPtrToVal(res.TotalUsers)), Valid: res.TotalUsers != nil},
		ActiveHalfYear:    pgtype.Int4{Int32: int32(utils.IntPtrToVal(res.ActiveHalfyear)), Valid: res.ActiveHalfyear != nil},
		ActiveMonth:       pgtype.Int4{Int32: int32(utils.IntPtrToVal(res.ActiveMonth)), Valid: res.ActiveMonth != nil},
		LocalPosts:        pgtype.Int4{Int32: int32(utils.IntPtrToVal(res.LocalPosts)), Valid: res.LocalPosts != nil},
		LocalComments:     pgtype.Int4{Int32: int32(utils.IntPtrToVal(res.LocalComments)), Valid: res.LocalComments != nil},
	})
	if err != nil {
		return err
	}

	// set the latest crawl_id on the server
	err = b.queries.UpdateServerLastCrawlID(ctx, db.UpdateServerLastCrawlIDParams{
		ID:          s.ID,
		LastCrawlID: c.ID,
	})
	if err != nil {
		return err
	}

	// start a transaction to add the peers
	tx, err := b.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck
	qtx := b.queries.WithTx(tx)

	// add the peers if they are not already in the db
	err = qtx.CreateServersFromDomainList(ctx, res.Peers)
	if err != nil {
		return err
	}

	// update the relations between the server and the peers
	err = qtx.UpdatePeeringRelationships(ctx, db.UpdatePeeringRelationshipsParams{
		ServerID: s.ID,
		Domains:  res.Peers,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
