package business

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/cyclimse/fediverse-blahaj/internal/db"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/cyclimse/fediverse-blahaj/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (cached *cachedErrorCodeDescriptions) Description(ctx context.Context, queries *db.Queries, code models.CrawlErrCode) (string, error) {
	if cached.LastUpdate.Before(time.Now().Add(-cached.ExpireAfter)) {
		rows, err := queries.ListErrorCodeDescriptions(ctx)
		if err != nil {
			return "", err
		}

		cached.LastUpdate = time.Now()
		for _, row := range rows {
			cached.Descriptions[models.CrawlErrCode(row.ErrorCode)] = row.Description
		}
	}

	desc, ok := cached.Descriptions[code]
	if !ok {
		// ultimately, it's not a big deal if we don't have a description
		slog.ErrorContext(ctx, "failed to get error code description", "code", code)
		return "", nil
	}
	return desc, nil
}

func (b *Business) GetInstanceByID(ctx context.Context, id uuid.UUID) (models.FediverseInstance, error) {
	row, err := b.queries.GetInstanceWithLastCrawlByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.FediverseInstance{}, ErrInstanceNotFound
		}
		return models.FediverseInstance{}, err
	}

	instance := models.FediverseInstance{
		ID:     row.ID.Bytes,
		Domain: row.Domain,
		Status: string(row.Status),

		SoftwareName: utils.ValToPtr(row.SoftwareName.String, row.SoftwareName.Valid),

		LastCrawl: &models.Crawl{
			ID:     row.LastCrawlID.Bytes,
			Domain: row.Domain,

			Status: models.CrawlStatus(row.Status_2),

			StartedAt:  row.StartedAt.Time,
			FinishedAt: row.FinishedAt.Time,

			Peers:         nil,
			NumberOfPeers: utils.ValToPtr(row.NumberOfPeers.Int32, row.NumberOfPeers.Valid),

			SoftwareName:    utils.ValToPtr(row.SoftwareName_2.String, row.SoftwareName_2.Valid),
			SoftwareVersion: utils.ValToPtr(row.SoftwareVersion.String, row.SoftwareVersion.Valid),

			OpenRegistrations: utils.ValToPtr(row.OpenRegistrations.Bool, row.OpenRegistrations.Valid),
			TotalUsers:        utils.ValToPtr(row.TotalUsers.Int32, row.TotalUsers.Valid),
			ActiveHalfyear:    utils.ValToPtr(row.ActiveHalfYear.Int32, row.ActiveHalfYear.Valid),
			ActiveMonth:       utils.ValToPtr(row.ActiveMonth.Int32, row.ActiveMonth.Valid),
			LocalPosts:        utils.ValToPtr(row.LocalPosts.Int32, row.LocalPosts.Valid),
			LocalComments:     utils.ValToPtr(row.LocalComments.Int32, row.LocalComments.Valid),

			RawNodeinfo: json.RawMessage(row.RawNodeinfo),
			Addresses:   row.Addresses,
		},
	}

	if row.ErrorMsg.Valid {
		d, err := b.errorCodeDescriptions.Description(ctx, b.queries, models.CrawlErrCode(row.ErrorCode.CrawlErrorCode))
		if err != nil {
			return models.FediverseInstance{}, err
		}

		instance.LastCrawl.Err = &models.CrawlError{
			Msg:         row.ErrorMsg.String,
			Code:        models.CrawlErrCode(row.ErrorCode.CrawlErrorCode),
			Description: d,
		}
	}

	return instance, nil
}

func (b *Business) ListInstances(ctx context.Context, page, pageSize int32) ([]models.FediverseInstance, int64, error) {
	rows, err := b.queries.ListInstancesPaginated(ctx, db.ListInstancesPaginatedParams{
		Limit:      pageSize,
		Offset:     (page - 1) * pageSize,
		TotalUsers: pgtype.Int4{Int32: smallServerThreshold, Valid: true},
	})
	if err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return nil, 0, nil
	}
	total := rows[0].TotalCount

	instances := make([]models.FediverseInstance, 0, len(rows))
	for _, row := range rows {
		instance := models.FediverseInstance{
			ID:     row.ID.Bytes,
			Domain: row.Domain,
			Status: string(row.Status),

			SoftwareName: utils.ValToPtr(row.SoftwareName.String, row.SoftwareName.Valid),

			LastCrawl: &models.Crawl{
				ID:     row.LastCrawlID.Bytes,
				Domain: row.Domain,

				Status: models.CrawlStatus(row.Status_2),

				StartedAt:  row.StartedAt.Time,
				FinishedAt: row.FinishedAt.Time,

				Peers:         nil,
				NumberOfPeers: utils.ValToPtr(row.NumberOfPeers.Int32, row.NumberOfPeers.Valid),

				SoftwareName:    utils.ValToPtr(row.SoftwareName.String, row.SoftwareName.Valid),
				SoftwareVersion: utils.ValToPtr(row.SoftwareVersion.String, row.SoftwareVersion.Valid),

				OpenRegistrations: utils.ValToPtr(row.OpenRegistrations.Bool, row.OpenRegistrations.Valid),
				TotalUsers:        utils.ValToPtr(row.TotalUsers.Int32, row.TotalUsers.Valid),
				ActiveHalfyear:    utils.ValToPtr(row.ActiveHalfYear.Int32, row.ActiveHalfYear.Valid),
				ActiveMonth:       utils.ValToPtr(row.ActiveMonth.Int32, row.ActiveMonth.Valid),
				LocalPosts:        utils.ValToPtr(row.LocalPosts.Int32, row.LocalPosts.Valid),
				LocalComments:     utils.ValToPtr(row.LocalComments.Int32, row.LocalComments.Valid),

				RawNodeinfo: json.RawMessage(row.RawNodeinfo),
				Addresses:   row.Addresses,
			},
		}

		if row.ErrorMsg.Valid {
			d, err := b.errorCodeDescriptions.Description(ctx, b.queries, models.CrawlErrCode(row.ErrorCode.CrawlErrorCode))
			if err != nil {
				return nil, 0, err
			}

			instance.LastCrawl.Err = &models.CrawlError{
				Msg:         row.ErrorMsg.String,
				Code:        models.CrawlErrCode(row.ErrorCode.CrawlErrorCode),
				Description: d,
			}
		}

		instances = append(instances, instance)
	}

	return instances, total, nil
}

func (b *Business) ListCrawlsForInstance(ctx context.Context, instanceID uuid.UUID, page, pageSize int32) ([]models.Crawl, int64, error) {
	rows, err := b.queries.ListCrawlsPaginated(ctx, db.ListCrawlsPaginatedParams{
		InstanceID: pgtype.UUID{Bytes: instanceID, Valid: true},
		Limit:      pageSize,
		Offset:     (page - 1) * pageSize,
	})
	if err != nil {
		return nil, 0, err
	}

	if len(rows) == 0 {
		return nil, 0, nil
	}
	total := rows[0].TotalCount

	crawls := make([]models.Crawl, 0, len(rows))

	for _, row := range rows {
		c := models.Crawl{
			ID: row.ID.Bytes,

			Status: models.CrawlStatus(row.Status),

			StartedAt:  row.StartedAt.Time,
			FinishedAt: row.FinishedAt.Time,

			Peers:         nil,
			NumberOfPeers: utils.ValToPtr(row.NumberOfPeers.Int32, row.NumberOfPeers.Valid),

			SoftwareName:    utils.ValToPtr(row.SoftwareName.String, row.SoftwareName.Valid),
			SoftwareVersion: utils.ValToPtr(row.SoftwareVersion.String, row.SoftwareVersion.Valid),

			OpenRegistrations: utils.ValToPtr(row.OpenRegistrations.Bool, row.OpenRegistrations.Valid),
			TotalUsers:        utils.ValToPtr(row.TotalUsers.Int32, row.TotalUsers.Valid),
			ActiveHalfyear:    utils.ValToPtr(row.ActiveHalfYear.Int32, row.ActiveHalfYear.Valid),
			ActiveMonth:       utils.ValToPtr(row.ActiveMonth.Int32, row.ActiveMonth.Valid),
			LocalPosts:        utils.ValToPtr(row.LocalPosts.Int32, row.LocalPosts.Valid),
			LocalComments:     utils.ValToPtr(row.LocalComments.Int32, row.LocalComments.Valid),

			RawNodeinfo: json.RawMessage(row.RawNodeinfo),
			Addresses:   row.Addresses,
		}

		if row.ErrorMsg.Valid {
			d, err := b.errorCodeDescriptions.Description(ctx, b.queries, models.CrawlErrCode(row.ErrorCode.CrawlErrorCode))
			if err != nil {
				return nil, 0, err
			}

			c.Err = &models.CrawlError{
				Msg:         row.ErrorMsg.String,
				Code:        models.CrawlErrCode(row.ErrorCode.CrawlErrorCode),
				Description: d,
			}
		}

		crawls = append(crawls, c)
	}

	return crawls, total, nil
}
