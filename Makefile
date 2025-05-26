include .env
BIN = cf-search

all:
	@templ generate
	@go build -o bin/$(BIN) cmd/web/*

cli-build:
	@go build -o bin/$(CLIBIN) cmd/cli/*

run: all
	@./bin/$(BIN)

sqlite:
	@sqlite3 $(DB)

migrate:
	migrate -database "sqlite3://$(DB)" -path ./migrations up
