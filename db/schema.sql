CREATE TYPE server_status AS ENUM (
  'active',
  'inactive',
  'deleted',
  'error',
  'unknown'
);


CREATE TABLE servers (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  domain varchar(255) NOT NULL UNIQUE,
  status server_status NOT NULL DEFAULT 'active',
  created_at timestamptz NOT NULL DEFAULT NOW(),
  deleted_at timestamptz,
  updated_at timestamptz,
  software_name varchar(255),
  -- can be null if we haven't crawled it yet
  last_crawl_id uuid
);


CREATE TABLE crawls (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  server_id uuid REFERENCES servers(id) NOT NULL,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  number_of_peers integer NOT NULL CHECK (number_of_peers >= 0),
  open_registrations boolean NOT NULL,
  total_users integer CHECK (total_users >= 0),
  active_half_year integer CHECK (active_half_year >= 0),
  active_month integer CHECK (active_month >= 0),
  local_posts integer CHECK (local_posts >= 0),
  local_comments integer CHECK (local_comments >= 0)
);


ALTER TABLE servers
ADD CONSTRAINT last_crawl_id FOREIGN KEY (last_crawl_id) REFERENCES crawls(id);


CREATE TABLE peering_relationships (
  server_id uuid REFERENCES servers(id),
  peer_id uuid REFERENCES servers(id),
  PRIMARY KEY (server_id, peer_id)
);


CREATE INDEX peering_relationships_server_id_idx ON peering_relationships (server_id);