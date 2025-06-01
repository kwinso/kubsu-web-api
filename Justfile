set dotenv-load := true

migrate-up:
  migrate -path sql/migrations -database $DATABASE_URL up

migrate-down:
  migrate -path sql/migrations -database $DATABASE_URL down

db-gen:
  sqlc generate

dev:
  gow run cmd/server/main.go -e=go,html