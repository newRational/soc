ifeq ($(POSTGRES_SETUP_PROD),)
	POSTGRES_SETUP_PROD := user=user password=pass dbname=prod host=localhost port=8002 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/
MOCKGEN_TAG=1.6.0
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/migrations

.PHONY: .migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: .migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_PROD)" up

.PHONY: .migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_PROD)" down

.PHONY: .dc-up
dc-up:
	sudo docker-compose up -d postgres zookeeper kafka1 kafka2 kafka3 redis

.PHONY: .dc-stop
dc-stop:
	sudo docker-compose stop

.PHONY: .dc-down
dc-down:
	sudo docker-compose down