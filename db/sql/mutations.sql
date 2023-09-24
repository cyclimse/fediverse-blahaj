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