CREATE TYPE server_status AS ENUM ('unknown', 'online', 'offline');


CREATE TYPE crawl_status AS ENUM (
  'unknown',
  -- crawling was successful
  'completed',
  -- generic failure not related to the crawler itself
  'failed',
  -- crawler was blocked by the server (e.g. robots.txt)
  'blocked',
  -- crawler timed out
  'timeout',
  -- error in the crawler itself
  'internal_error'
);


CREATE TABLE servers (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  domain varchar(255) NOT NULL UNIQUE,
  status server_status NOT NULL DEFAULT 'unknown',
  created_at timestamptz NOT NULL DEFAULT NOW(),
  deleted_at timestamptz,
  updated_at timestamptz,
  software_name varchar(255),
  -- can be null if we haven't crawled it yet
  last_crawl_id uuid
);


CREATE INDEX servers_domain_idx ON servers (domain);


CREATE TABLE crawls (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  server_id uuid REFERENCES servers(id) NOT NULL,
  status crawl_status NOT NULL DEFAULT 'unknown',
  error_msg text DEFAULT NULL,
  started_at timestamptz NOT NULL DEFAULT NOW(),
  -- from nodeinfo
  software_name varchar(255),
  number_of_peers integer CHECK (number_of_peers >= 0),
  open_registrations boolean,
  total_users integer CHECK (total_users >= 0),
  active_half_year integer CHECK (active_half_year >= 0),
  active_month integer CHECK (active_month >= 0),
  local_posts integer CHECK (local_posts >= 0),
  local_comments integer CHECK (local_comments >= 0)
);


CREATE INDEX crawls_server_id_idx ON crawls (server_id);


ALTER TABLE servers
ADD CONSTRAINT last_crawl_id FOREIGN KEY (last_crawl_id) REFERENCES crawls(id);


CREATE TABLE peering_relationships (
  server_id uuid REFERENCES servers(id),
  peer_id uuid REFERENCES servers(id),
  PRIMARY KEY (server_id, peer_id)
);


CREATE INDEX peering_relationships_server_id_idx ON peering_relationships (server_id);