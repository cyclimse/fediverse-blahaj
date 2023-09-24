CREATE TYPE instance_status AS ENUM ('unknown', 'up', 'down', 'unhealthy');


CREATE TYPE crawl_status AS ENUM ('unknown', 'completed', 'failed');


CREATE TYPE crawl_error_code AS ENUM (
  'unknown',
  'timeout',
  'domain_not_found',
  'unreachable',
  'invalid_nodeinfo',
  'nodeinfo_version_not_supported_by_crawler',
  'invalid_json',
  'blocked_by_robots_txt',
  'software_not_supported_by_crawler',
  'software_version_not_supported_by_crawler',
  'internal_error'
);


CREATE TABLE instance (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  domain varchar(512) UNIQUE NOT NULL,
  status instance_status NOT NULL DEFAULT 'unknown',
  created_at timestamptz NOT NULL DEFAULT NOW(),
  deleted_at timestamptz,
  updated_at timestamptz,
  -- from nodeinfo
  software_name varchar(255),
  -- can be null if we haven' t crawled it yet last_crawl_id uuid
  last_crawl_id uuid
);


CREATE INDEX instance_domain_idx ON instance (domain);


CREATE TABLE crawl (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  instance_id uuid REFERENCES instance(id) NOT NULL,
  status crawl_status NOT NULL DEFAULT 'unknown',
  -- if the crawl failed, this is the reason why
  -- error_code can be returned to the user, error_msg is for debugging
  error_code crawl_error_code DEFAULT NULL,
  error_msg varchar(1024) DEFAULT NULL,
  -- when the crawl was started and finished
  -- should not be null as we only write to the table once the crawl is finished
  started_at timestamptz NOT NULL,
  finished_at timestamptz NOT NULL,
  -- from nodeinfo
  software_name varchar(255),
  software_version varchar(255),
  number_of_peers integer CHECK (number_of_peers >= 0),
  open_registrations boolean,
  total_users integer CHECK (total_users >= 0),
  active_half_year integer CHECK (active_half_year >= 0),
  active_month integer CHECK (active_month >= 0),
  local_posts integer CHECK (local_posts >= 0),
  local_comments integer CHECK (local_comments >= 0),
  -- jsonb object containing the raw nodeinfo response
  -- useful for backfilling data
  raw_nodeinfo jsonb,
  -- ip address of the instance. not displayed publicly, but useful for
  -- debugging and blocking.
  addresses inet []
);


CREATE INDEX crawl_instance_id_idx ON crawl (instance_id);


ALTER TABLE instance
ADD CONSTRAINT last_crawl_id FOREIGN KEY (last_crawl_id) REFERENCES crawl(id);


CREATE TABLE peering_relationship (
  instance_id uuid REFERENCES instance(id),
  peer_id uuid REFERENCES instance(id),
  PRIMARY KEY (instance_id, peer_id)
);


CREATE INDEX peering_relationship_instance_id_idx ON peering_relationship (instance_id);


CREATE TABLE crawl_errors (
  error_code crawl_error_code PRIMARY KEY,
  description varchar(1024) NOT NULL
);