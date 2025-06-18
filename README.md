# go-sql-practice

## Dev Setup

- copy `env.example` to `.env`
- install goose
- install sqlc
- goose up/down
- sqlc generate

---

## Deploy

Render
- Build: `go install github.com/pressly/goose/v3/cmd/goose@latest && goose up && go build -tags netgo -ldflags '-s -w' -o app`
- Start Command: `./app`