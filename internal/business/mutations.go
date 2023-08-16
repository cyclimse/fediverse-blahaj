package business

import (
	"context"
	"fmt"

	"github.com/cyclimse/fediverse-blahaj/internal/db"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/cyclimse/fediverse-blahaj/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func (b *Business) AddServer(ctx context.Context, server models.FediverseServer) (db.Server, error) {
	if server.Domain == "" {
		return db.Server{}, fmt.Errorf("domain is empty")
	}

	s, error := b.queries.CreateServer(ctx, db.CreateServerParams{
		Domain:       server.Domain,
		SoftwareName: pgtype.Text{String: utils.StringPtrToVal(server.SoftwareName), Valid: server.SoftwareName != nil},
	})
	if error != nil {
		return db.Server{}, error
	}
	return s, nil
}

func (b *Business) AddFailedCrawlToServer(ctx context.Context, res models.FediverseServer, s db.Server) error {
	if res.CrawlErr == nil {
		return fmt.Errorf("crawl result has no error")
	}

	c, err := b.queries.CreateFailedCrawl(ctx, db.CreateFailedCrawlParams{
		ServerID: s.ID,
		Status:   db.CrawlStatus(res.CrawlStatus),
		// res.CrawlErr is not nil, so we can safely use it
		ErrorMsg: pgtype.Text{String: res.CrawlErr.Error(), Valid: true},
	})
	if err != nil {
		return err
	}

	// set the latest crawl_id on the server
	err = b.queries.UpdateServerLastCrawlID(ctx, db.UpdateServerLastCrawlIDParams{
		LastCrawlID: c.ID,
		// TODO: maybe we should not set it to offline if it is a temporary error
		// do this asynchroniously in a cronjob
		Status: db.ServerStatusOffline,
		ID:     s.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (b *Business) AddCompletedCrawlToServer(ctx context.Context, res models.FediverseServer, s db.Server) error {
	c, err := b.queries.CreateCompletedCrawl(ctx, db.CreateCompletedCrawlParams{
		ServerID:          s.ID,
		NumberOfPeers:     pgtype.Int4{Int32: int32(utils.IntPtrToVal(res.NumberOfPeers)), Valid: res.NumberOfPeers != nil},
		OpenRegistrations: pgtype.Bool{Bool: utils.BoolPtrToVal(res.OpenRegistrations), Valid: res.OpenRegistrations != nil},
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
		Status:      db.ServerStatusOnline,
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
