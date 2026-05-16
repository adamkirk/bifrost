
DC=docker compose -f ./docker-compose.yml -p bifrost

.PHONY: npm \
	up down restart dc build \
	setup-hooks

npm:
	@[ $$(node -v | tr -d v | cut -d. -f1) -ge 25 ] || { echo "Error: node $$(node -v) < required v25"; exit 1; }
	@[ $$(npm -v | cut -d. -f1) -ge 11 ] || { echo "Error: npm $$(npm -v) < required v11"; exit 1; }
	npm install

up:
	$(DC) up -d

down:
	$(DC) down --remove-orphans

restart: down up

build:
	$(DC) build

dc:
	@echo "$(DC)"

setup-hooks:
	git config core.hooksPath .githooks

exec-%:
	$(DC) exec -it $* bash

tail-%:
	$(DC) logs -f $*