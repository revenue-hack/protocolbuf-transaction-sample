.PHONY: create dryrun migrate/prd migrate/dev migrate/test migrate/local rollback/local status

OS := $(shell uname -s)
ifeq ($(OS),Darwin)
	export SHELL := $(shell echo $$SHELL)
endif

## create
create:
	sql-migrate new "${file}"
## dryrun
dryrun:
	sql-migrate up --env=local -dryrun
## migrate
migrate/local:
	sql-migrate up --env="local"
## rollback
rollback/local:
	sql-migrate down --env="local"
## status
status:
	sql-migrate status


