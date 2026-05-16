package config

import "go-sql-practice/internal/database"

type ApiConfig struct {
	Db *database.Queries
}
