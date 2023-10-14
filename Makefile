run-http:
	go run ./cmd/to-do-list-http

build-http:
	go build -o bin/to-do-list-http ./cmd/to-do-list-http

mod-vendor:
	go mod vendor

docker-compose:
	docker-compose up -d

init-app:
	mod-vendor && docker-compose