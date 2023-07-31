package business

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/cyclimse/fediverse-blahaj/internal/db"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/cyclimse/fediverse-blahaj/internal/utils"
	"github.com/google/uuid"
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

func (b *Business) GetServerByID(ctx context.Context, id uuid.UUID) (*models.FediverseServer, error) {
	row, err := b.queries.GetServerWithLastCrawlByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrServerNotFound
		}
		return nil, err
	}
	return &models.FediverseServer{
		ID:     row.ID.Bytes,
		Domain: row.Domain,

		Peers:         nil,
		NumberOfPeers: row.NumberOfPeers,

		SoftwareName: row.SoftwareName.String,

		OpenRegistrations: row.OpenRegistrations,
		TotalUsers:        utils.IntValToPtr(row.TotalUsers.Int32, row.TotalUsers.Valid),
		ActiveHalfyear:    utils.IntValToPtr(row.ActiveHalfYear.Int32, row.ActiveHalfYear.Valid),
		ActiveMonth:       utils.IntValToPtr(row.ActiveMonth.Int32, row.ActiveMonth.Valid),
		LocalPosts:        utils.IntValToPtr(row.LocalPosts.Int32, row.LocalPosts.Valid),
		LocalComments:     utils.IntValToPtr(row.LocalComments.Int32, row.LocalComments.Valid),
	}, nil
}

func (b *Business) ListServers(ctx context.Context, page, pageSize int32) ([]models.FediverseServer, error) {
	res, err := b.queries.ListSeversPaginated(ctx, db.ListSeversPaginatedParams{
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return nil, err
	}

	servers := make([]models.FediverseServer, len(res))
	for i := range res {
		servers[i] = models.FediverseServer{
			ID:     res[i].ID.Bytes,
			Domain: res[i].Domain,

			Peers:         nil,
			NumberOfPeers: res[i].NumberOfPeers,

			SoftwareName: res[i].SoftwareName.String,

			OpenRegistrations: res[i].OpenRegistrations,
			TotalUsers:        utils.IntValToPtr(res[i].TotalUsers.Int32, res[i].TotalUsers.Valid),
			ActiveHalfyear:    utils.IntValToPtr(res[i].ActiveHalfYear.Int32, res[i].ActiveHalfYear.Valid),
			ActiveMonth:       utils.IntValToPtr(res[i].ActiveMonth.Int32, res[i].ActiveMonth.Valid),
			LocalPosts:        utils.IntValToPtr(res[i].LocalPosts.Int32, res[i].LocalPosts.Valid),
			LocalComments:     utils.IntValToPtr(res[i].LocalComments.Int32, res[i].LocalComments.Valid),
		}
	}

	return servers, nil
}
