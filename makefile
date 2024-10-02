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


BASE_URL = http://localhost:8080/api/v1
add_song:
	curl -X POST "${BASE_URL}/songs" -H "Content-Type: application/json" --data '${data}' | jq

delete_song:
	curl -X DELETE "${BASE_URL}/songs/${id}" | jq

get_song_text:
	curl -X GET "${BASE_URL}/songs/${id}/lyrics${query}" | jq

get_songs:
	curl -X GET "${BASE_URL}/songs${query}" | jq

get_song_details:
	curl -X GET "${BASE_URL}/info${query}" | jq

update_song:
	curl -X PATCH "${BASE_URL}/songs/${id}" -H "Content-Type: application/json" --data '${data}'


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