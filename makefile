MIGRATE = migrate
MIGR_DIR = migrations
DB_SOURCE = postgres://amicie:admin@localhost:5432/$(1)?sslmode=disable
MIGRATE_BODY = ${MIGRATE} -path ${MIGR_DIR} -database $(call DB_SOURCE,$(db))

migrate_new: check_seq
	@${MIGRATE} create -ext sql -dir migrations/ -seq ${seq}

migrate: check_db
	@${MIGRATE_BODY} up

migrate_force: check_db check_v
	@${MIGRATE_BODY} force ${v}

migrate_down:
	@${MIGRATE_BODY} down

check_db:
ifndef db
	@$(error parameter db is required [database name])
endif

check_v:
ifndef v
	@$(error parameter v is required [version])
endif

check_seq:
ifndef seq
	@$(error parameter seq is required [sequence])
endif