set dotenv-load := true

migrate-up:
  migrate -path db/migrations -database $DATABASE_URL up

migrate-down:
  migrate -path db/migrations -database $DATABASE_URL down

db-gen:
  sqlc generate