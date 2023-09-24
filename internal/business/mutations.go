package business

import (
	"context"
	"fmt"

	"github.com/cyclimse/fediverse-blahaj/internal/db"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/cyclimse/fediverse-blahaj/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

func (b *Business) AddInstance(ctx context.Context, domain string, software *string) (db.Instance, error) {
	if domain == "" {
		return db.Instance{}, fmt.Errorf("domain is empty")
	}

	s, err := b.queries.CreateInstance(ctx, db.CreateInstanceParams{
		Domain:       domain,
		SoftwareName: pgtype.Text{String: utils.StringPtrToVal(software), Valid: software != nil},
	})
	if err != nil {
		return db.Instance{}, err
	}
	return s, nil
}

func (b *Business) AddCrawlToFediverseInstance(ctx context.Context, crawl models.Crawl, instance db.Instance) error {
	params := db.CreateCrawlParams{
		InstanceID: instance.ID,

		Status: db.CrawlStatus(crawl.Status),

		StartedAt:  pgtype.Timestamptz{Time: crawl.StartedAt, Valid: true},
		FinishedAt: pgtype.Timestamptz{Time: crawl.FinishedAt, Valid: true},

		SoftwareName:    pgtype.Text{String: utils.StringPtrToVal(crawl.SoftwareName), Valid: crawl.SoftwareName != nil},
		SoftwareVersion: pgtype.Text{String: utils.StringPtrToVal(crawl.SoftwareVersion), Valid: crawl.SoftwareVersion != nil},

		NumberOfPeers:     pgtype.Int4{Int32: int32(utils.IntPtrToVal(crawl.NumberOfPeers)), Valid: crawl.NumberOfPeers != nil},
		OpenRegistrations: pgtype.Bool{Bool: utils.BoolPtrToVal(crawl.OpenRegistrations), Valid: crawl.OpenRegistrations != nil},
		TotalUsers:        pgtype.Int4{Int32: int32(utils.IntPtrToVal(crawl.TotalUsers)), Valid: crawl.TotalUsers != nil},
		ActiveHalfYear:    pgtype.Int4{Int32: int32(utils.IntPtrToVal(crawl.ActiveHalfyear)), Valid: crawl.ActiveHalfyear != nil},
		ActiveMonth:       pgtype.Int4{Int32: int32(utils.IntPtrToVal(crawl.ActiveMonth)), Valid: crawl.ActiveMonth != nil},
		LocalPosts:        pgtype.Int4{Int32: int32(utils.IntPtrToVal(crawl.LocalPosts)), Valid: crawl.LocalPosts != nil},
		LocalComments:     pgtype.Int4{Int32: int32(utils.IntPtrToVal(crawl.LocalComments)), Valid: crawl.LocalComments != nil},

		RawNodeinfo: []byte(crawl.RawNodeinfo),
		Addresses:   crawl.Addresses,
	}

	instanceStatus := db.InstanceStatusUp
	if crawl.Err != nil {
		params.ErrorMsg = pgtype.Text{String: crawl.Err.Error(), Valid: true}
		params.ErrorCode = db.NullCrawlErrorCode{CrawlErrorCode: db.CrawlErrorCode(crawl.Err.Code), Valid: true}
		instanceStatus = db.InstanceStatusDown
	}

	c, err := b.queries.CreateCrawl(ctx, params)
	if err != nil {
		return err
	}

	// set the latest crawl_id on the server
	err = b.queries.UpdateInstanceFromLastCrawl(ctx, db.UpdateInstanceFromLastCrawlParams{
		ID:     instance.ID,
		Status: instanceStatus,

		// While this is considered immutable, it is not enforced by the db.
		// It's important to update it here because when creating instances in bulk,
		// the software name is not known yet and therefore not set.
		SoftwareName: pgtype.Text{String: utils.StringPtrToVal(crawl.SoftwareName), Valid: crawl.SoftwareName != nil},

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
	// TODO: filter blacklisted domains
	err = qtx.CreateInstancesFromDomainList(ctx, crawl.Peers)
	if err != nil {
		return err
	}

	// update the relations between the server and the peers
	err = qtx.UpdatePeeringRelationships(ctx, db.UpdatePeeringRelationshipsParams{
		InstanceID: instance.ID,
		Domains:    crawl.Peers,
	})
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
