-- name: GetInstanceByDomain :one
SELECT *
FROM instance
WHERE domain = $1
LIMIT 1;


-- name: GetInstanceWithLastCrawlByID :one
SELECT *
FROM instance
  JOIN crawl ON crawl.id = instance.last_crawl_id
WHERE instance.id = $1
  AND instance.deleted_at IS NULL
LIMIT 1;


-- name: GetPeersIDsByInstanceID :many
SELECT peer_id
FROM peering_relationship
  JOIN instance ON instance.id = peer_id
  AND instance.deleted_at IS NULL
WHERE instance_id = $1;


-- TODO: these types of paginated queries are not efficient
--       we should use a cursor instead or a CTE
-- name: ListInstancesPaginated :many
SELECT *,
  COUNT(*) OVER() AS total_count
FROM instance
  JOIN crawl ON crawl.id = instance.last_crawl_id
WHERE deleted_at IS NULL
  AND total_users > $3
ORDER BY total_users DESC
LIMIT $1 OFFSET $2;


-- name: ListCrawlsPaginated :many
SELECT *,
  COUNT(*) OVER() AS total_count
FROM crawl
WHERE instance_id = $1
ORDER BY started_at DESC
LIMIT $2 OFFSET $3;


-- name: ListErrorCodeDescriptions :many
SELECT *
FROM crawl_errors;


-- name: CreateInstance :one
INSERT INTO instance (domain, software_name)
VALUES ($1, $2)
RETURNING *;


-- name: CreateInstancesFromDomainList :exec
INSERT INTO instance (domain)
SELECT domain
FROM unnest(@domains::varchar(255) []) domain ON CONFLICT DO NOTHING;


-- name: CreateCrawl :one
INSERT INTO crawl (
    instance_id,
    status,
    error_code,
    error_msg,
    started_at,
    finished_at,
    software_name,
    software_version,
    number_of_peers,
    open_registrations,
    total_users,
    active_half_year,
    active_month,
    local_posts,
    local_comments,
    raw_nodeinfo,
    addresses
  )
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15,
    $16,
    $17
  )
RETURNING *;


-- name: UpdatePeeringRelationships :exec
INSERT INTO peering_relationship (instance_id, peer_id)
SELECT $1,
  id
FROM instance
WHERE domain = ANY(@domains::varchar(255) []) ON CONFLICT DO NOTHING;


-- name: UpdateInstanceFromLastCrawl :exec
UPDATE instance
SET last_crawl_id = $2,
  status = $3,
  software_name = $4,
  updated_at = NOW()
WHERE id = $1;


-- name: DeleteInstanceByID :exec
DELETE FROM instance
WHERE id = $1;