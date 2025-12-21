

start-local:
	docker compose up -d
	go run ./cmd/main.go run rest

stop-local:
	docker compose down