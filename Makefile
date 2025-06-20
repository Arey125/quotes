include .env
BIN = quotes

.PHONY: all tailwind run sqlite migrate migrate-down

all: tailwind
	@templ generate
	@go build -o bin/$(BIN) cmd/web/*

tailwind:
	@$(TAILWIND) -i app.css -o ./static/tailwind-output.css --minify

run: all
	@./bin/$(BIN)

sqlite:
	@sqlite3 $(DB)

migrate:
	migrate -database "sqlite3://$(DB)" -path ./migrations up

migrate-down:
	migrate -database "sqlite3://$(DB)" -path ./migrations down 1
