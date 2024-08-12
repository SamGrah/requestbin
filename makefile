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

.PHONY: dev
dev:
	templ generate && go build -o tmp/app ./cmd/main.go && air -c .air.toml

.PHONY: build
build:
	make tailwind-build
	make templ-generate
	go build -ldflags "-X main.environment=production" -o ./bin/app .