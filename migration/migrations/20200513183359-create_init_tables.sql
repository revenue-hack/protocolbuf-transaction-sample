
-- +migrate Up

-- -----------------------------------------------------
-- Table users
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(40) NOT NULL,
  name VARCHAR(64) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  updated_at timestamp  NOT NULL DEFAULT CURRENT_TIMESTAMP(),
  PRIMARY KEY (`id`))ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXISTS users;

