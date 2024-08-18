.PHONY: tailwind-watch
tailwind-watch:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	./tailwindcss -i ./static/css/input.css -o ./static/css/style.min.css --minify

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: templ-watch
templ-watch:
	templ generate --watch 

.PHONY: create-db
db-init:
ifeq ("$(wildcard database.db)","")
	sqlite3 database.db < db-schema.sql
else
	@:
endif
	
.PHONY: dev
dev:
	make db-init
	templ generate
	go build -o tmp/app ./cmd/main.go
	air -c .air.toml

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -ldflags "-X main.environment=production" -o ./bin/app .