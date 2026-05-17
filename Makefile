
DC=docker compose --env-file ./.dockerenv -f ./docker-compose.yml -p bifrost
TTY_FLAGS=-it

ifeq ($(NO_TTY),1)
TTY_FLAGS=-T
endif

.PHONY: npm \
	prepare-env \
	up down restart dc build \
	hosts-txeh setup-certs-mkcert \
	setup-hooks \
	tail-% exec-% \
	api-tidy hook-api-fmtcheck

prepare-env:
	@[ -f .dockerenv ] || cp .dockerenv.example .dockerenv

npm:
	@[ $$(node -v | tr -d v | cut -d. -f1) -ge 25 ] || { echo "Error: node $$(node -v) < required v25"; exit 1; }
	@[ $$(npm -v | cut -d. -f1) -ge 11 ] || { echo "Error: npm $$(npm -v) < required v11"; exit 1; }
	npm install

up: prepare-env
	$(DC) up -d

down:
	$(DC) down --remove-orphans

restart: down up

build:
	$(DC) build

dc:
	@echo "$(DC)"

setup-certs-mkcert:
	@which mkcert > /dev/null || { echo "Error: mkcert not installed. Run: brew install mkcert"; exit 1; }
	@mkcert -install
	@mkdir -p .local/traefik/certs
	@cd .local/traefik/certs && mkcert "*.bifrost.test"

setup-hooks:
	git config core.hooksPath .githooks

exec-%:
	@$(DC) exec $(TTY_FLAGS) $* bash

tail-%:
	@$(DC) logs -f $*

api-tidy:
	@$(DC) exec $(TTY_FLAGS) api go mod tidy

api-test:
	@$(DC) exec $(TTY_FLAGS) api ./scripts/test.sh ./...

api-db-migrate:
	@$(DC) run $(TTY_FLAGS) migrate up

api-db-wipe:
	@$(DC) run $(TTY_FLAGS) migrate drop

api-sqlc-gen:
	$(DC) exec $(TTY_FLAGS) api sqlc generate

hook-api-fmtcheck:
	$(DC) exec $(TTY_FLAGS) api ./scripts/fmt_check.sh

hosts-txeh:
	txeh add 127.0.0.1 \
		bifrost.test api.bifrost.test \
		traefik.bifrost.test pgadmin.bifrost.test \
		--comment "bifrost-dev"