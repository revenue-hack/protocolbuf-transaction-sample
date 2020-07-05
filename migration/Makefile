.PHONY: migrate/local rollback/local

OS := $(shell uname -s)
ifeq ($(OS),Darwin)
	export SHELL := $(shell echo $$SHELL)
endif

## migrate
migrate/local:
	sql-migrate up --env="local"
## rollback
rollback/local:
	sql-migrate down --env="local"

