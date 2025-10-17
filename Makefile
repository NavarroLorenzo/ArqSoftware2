SHELL := /bin/bash

.PHONY: tidy build up down logs

tidy:
	go work sync || true
	go mod tidy -C users-api
	go mod tidy -C canchas-api
	go mod tidy -C reservas-api
	go mod tidy -C search-api

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down -v

logs:
	docker compose logs -f --tail=100
