package business

import (
	"context"

	"github.com/cyclimse/fediverse-blahaj/internal/db"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/cyclimse/fediverse-blahaj/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

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
		NumberOfPeers: utils.ValToPtr(row.NumberOfPeers.Int32, row.NumberOfPeers.Valid),

		SoftwareName: utils.ValToPtr(row.SoftwareName.String, row.SoftwareName.Valid),

		OpenRegistrations: utils.ValToPtr(row.OpenRegistrations.Bool, row.OpenRegistrations.Valid),
		TotalUsers:        utils.ValToPtr(row.TotalUsers.Int32, row.TotalUsers.Valid),
		ActiveHalfyear:    utils.ValToPtr(row.ActiveHalfYear.Int32, row.ActiveHalfYear.Valid),
		ActiveMonth:       utils.ValToPtr(row.ActiveMonth.Int32, row.ActiveMonth.Valid),
		LocalPosts:        utils.ValToPtr(row.LocalPosts.Int32, row.LocalPosts.Valid),
		LocalComments:     utils.ValToPtr(row.LocalComments.Int32, row.LocalComments.Valid),
	}, nil
}

func (b *Business) ListServers(ctx context.Context, page, pageSize int32) ([]models.FediverseServer, int64, error) {
	res, err := b.queries.ListServersPaginated(ctx, db.ListServersPaginatedParams{
		Limit:      pageSize,
		Offset:     (page - 1) * pageSize,
		TotalUsers: pgtype.Int4{Int32: smallServerThreshold, Valid: true},
	})
	if err != nil {
		return nil, 0, err
	}

	if len(res) == 0 {
		return nil, 0, nil
	}
	total := res[0].TotalCount

	servers := make([]models.FediverseServer, 0, len(res))
	for i := range res {
		m := models.FediverseServer{
			ID:     res[i].ID.Bytes,
			Domain: res[i].Domain,

			Peers:         nil,
			NumberOfPeers: utils.ValToPtr(res[i].NumberOfPeers.Int32, res[i].NumberOfPeers.Valid),

			SoftwareName: utils.ValToPtr(res[i].SoftwareName.String, res[i].SoftwareName.Valid),

			OpenRegistrations: utils.ValToPtr(res[i].OpenRegistrations.Bool, res[i].OpenRegistrations.Valid),
			TotalUsers:        utils.ValToPtr(res[i].TotalUsers.Int32, res[i].TotalUsers.Valid),
			ActiveHalfyear:    utils.ValToPtr(res[i].ActiveHalfYear.Int32, res[i].ActiveHalfYear.Valid),
			ActiveMonth:       utils.ValToPtr(res[i].ActiveMonth.Int32, res[i].ActiveMonth.Valid),
			LocalPosts:        utils.ValToPtr(res[i].LocalPosts.Int32, res[i].LocalPosts.Valid),
			LocalComments:     utils.ValToPtr(res[i].LocalComments.Int32, res[i].LocalComments.Valid),
		}

		servers = append(servers, m)
	}

	return servers, total, nil
}
