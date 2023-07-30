// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createCrawl = `-- name: CreateCrawl :one
INSERT INTO crawls (
    server_id,
    number_of_peers,
    open_registrations,
    total_users,
    active_half_year,
    active_month,
    local_posts,
    local_comments
  )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, server_id, created_at, number_of_peers, open_registrations, total_users, active_half_year, active_month, local_posts, local_comments
`

type CreateCrawlParams struct {
	ServerID          pgtype.UUID
	NumberOfPeers     int32
	OpenRegistrations bool
	TotalUsers        pgtype.Int4
	ActiveHalfYear    pgtype.Int4
	ActiveMonth       pgtype.Int4
	LocalPosts        pgtype.Int4
	LocalComments     pgtype.Int4
}

func (q *Queries) CreateCrawl(ctx context.Context, arg CreateCrawlParams) (Crawl, error) {
	row := q.db.QueryRow(ctx, createCrawl,
		arg.ServerID,
		arg.NumberOfPeers,
		arg.OpenRegistrations,
		arg.TotalUsers,
		arg.ActiveHalfYear,
		arg.ActiveMonth,
		arg.LocalPosts,
		arg.LocalComments,
	)
	var i Crawl
	err := row.Scan(
		&i.ID,
		&i.ServerID,
		&i.CreatedAt,
		&i.NumberOfPeers,
		&i.OpenRegistrations,
		&i.TotalUsers,
		&i.ActiveHalfYear,
		&i.ActiveMonth,
		&i.LocalPosts,
		&i.LocalComments,
	)
	return i, err
}

const createServer = `-- name: CreateServer :one
INSERT INTO servers (domain, software_name)
VALUES ($1, $2)
RETURNING id, domain, status, created_at, deleted_at, updated_at, software_name, last_crawl_id
`

type CreateServerParams struct {
	Domain       string
	SoftwareName pgtype.Text
}

func (q *Queries) CreateServer(ctx context.Context, arg CreateServerParams) (Server, error) {
	row := q.db.QueryRow(ctx, createServer, arg.Domain, arg.SoftwareName)
	var i Server
	err := row.Scan(
		&i.ID,
		&i.Domain,
		&i.Status,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.UpdatedAt,
		&i.SoftwareName,
		&i.LastCrawlID,
	)
	return i, err
}

const createServersFromDomainList = `-- name: CreateServersFromDomainList :exec
INSERT INTO servers (domain)
SELECT domain
FROM unnest($1::varchar(255) []) domain ON CONFLICT DO NOTHING
`

func (q *Queries) CreateServersFromDomainList(ctx context.Context, domains []string) error {
	_, err := q.db.Exec(ctx, createServersFromDomainList, domains)
	return err
}

const deleteServerByID = `-- name: DeleteServerByID :exec
DELETE FROM servers
WHERE id = $1
`

func (q *Queries) DeleteServerByID(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteServerByID, id)
	return err
}

const getSeverByDomain = `-- name: GetSeverByDomain :one
SELECT id, domain, status, created_at, deleted_at, updated_at, software_name, last_crawl_id
FROM servers
WHERE domain = $1
LIMIT 1
`

func (q *Queries) GetSeverByDomain(ctx context.Context, domain string) (Server, error) {
	row := q.db.QueryRow(ctx, getSeverByDomain, domain)
	var i Server
	err := row.Scan(
		&i.ID,
		&i.Domain,
		&i.Status,
		&i.CreatedAt,
		&i.DeletedAt,
		&i.UpdatedAt,
		&i.SoftwareName,
		&i.LastCrawlID,
	)
	return i, err
}

const listSeversPaginated = `-- name: ListSeversPaginated :many
SELECT id, domain, status, created_at, deleted_at, updated_at, software_name, last_crawl_id
FROM servers
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListSeversPaginatedParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) ListSeversPaginated(ctx context.Context, arg ListSeversPaginatedParams) ([]Server, error) {
	rows, err := q.db.Query(ctx, listSeversPaginated, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Server
	for rows.Next() {
		var i Server
		if err := rows.Scan(
			&i.ID,
			&i.Domain,
			&i.Status,
			&i.CreatedAt,
			&i.DeletedAt,
			&i.UpdatedAt,
			&i.SoftwareName,
			&i.LastCrawlID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePeeringRelationships = `-- name: UpdatePeeringRelationships :exec
INSERT INTO peering_relationships (server_id, peer_id)
SELECT $1,
  id
FROM servers
WHERE domain = ANY($2::varchar(255) []) ON CONFLICT DO NOTHING
`

type UpdatePeeringRelationshipsParams struct {
	ServerID pgtype.UUID
	Domains  []string
}

func (q *Queries) UpdatePeeringRelationships(ctx context.Context, arg UpdatePeeringRelationshipsParams) error {
	_, err := q.db.Exec(ctx, updatePeeringRelationships, arg.ServerID, arg.Domains)
	return err
}

const updateServerLastCrawlID = `-- name: UpdateServerLastCrawlID :exec
UPDATE servers
SET last_crawl_id = $1
WHERE id = $2
`

type UpdateServerLastCrawlIDParams struct {
	LastCrawlID pgtype.UUID
	ID          pgtype.UUID
}

func (q *Queries) UpdateServerLastCrawlID(ctx context.Context, arg UpdateServerLastCrawlIDParams) error {
	_, err := q.db.Exec(ctx, updateServerLastCrawlID, arg.LastCrawlID, arg.ID)
	return err
}
