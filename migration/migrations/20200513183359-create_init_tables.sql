
-- +migrate Up

-- -----------------------------------------------------
-- Table users
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(40) NOT NULL,
  name VARCHAR(64) NOT NULL,
  created_at timestamp with time zone DEFAULT now() NOT NULL,
  updated_at timestamp with time zone DEFAULT now() NOT NULL,
  PRIMARY KEY (id),
  UNIQUE (name));

-- +migrate Down
DROP TABLE IF EXISTS users;

