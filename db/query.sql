-- name: GetSeverByDomain :one
SELECT *
FROM servers
WHERE domain = $1
LIMIT 1;


-- name: GetServerWithLastCrawlByID :one
SELECT *
FROM servers
  JOIN crawls ON crawls.id = servers.last_crawl_id
WHERE servers.id = $1
  AND servers.deleted_at IS NULL
LIMIT 1;


-- name: ListSeversPaginated :many
SELECT *
FROM servers
  JOIN crawls ON crawls.id = servers.last_crawl_id
WHERE servers.deleted_at IS NULL
ORDER BY servers.created_at DESC
LIMIT $1 OFFSET $2;


-- name: CreateServer :one
INSERT INTO servers (domain, software_name)
VALUES ($1, $2)
RETURNING *;


-- name: CreateServersFromDomainList :exec
INSERT INTO servers (domain)
SELECT domain
FROM unnest(@domains::varchar(255) []) domain ON CONFLICT DO NOTHING;


-- name: CreateCrawl :one
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
RETURNING *;


-- name: UpdatePeeringRelationships :exec
INSERT INTO peering_relationships (server_id, peer_id)
SELECT $1,
  id
FROM servers
WHERE domain = ANY(@domains::varchar(255) []) ON CONFLICT DO NOTHING;


-- name: UpdateServerLastCrawlID :exec
UPDATE servers
SET last_crawl_id = $1
WHERE id = $2;


-- name: DeleteServerByID :exec
DELETE FROM servers
WHERE id = $1;