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


-- name: GetCrawlerSeedDomains :many
-- Get all domains that have not been crawled yet. If there are none, return the domains with the oldest crawl.
SELECT domain
FROM instance
WHERE last_crawl_id IS NULL
  AND deleted_at IS NULL
UNION
(
  SELECT domain
  FROM instance
    JOIN crawl ON crawl.id = instance.last_crawl_id
  WHERE deleted_at IS NULL
  ORDER BY started_at ASC
)
LIMIT $1;