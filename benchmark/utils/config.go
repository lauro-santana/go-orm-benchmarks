package utils

import (
	"os"
)

// TODO: Add these ones to the .env file.
const (
	BulkInsertNumber     = 2000
	BatchSize            = 10000
	PageSize             = 10
	BulkInsertPageNumber = 100
	FindOneLoop          = 1
)

var PostgresDSN string

func init() {
	PostgresDSN = os.Getenv("POSTGRES_DSN")
	if PostgresDSN == "" {
		panic("POSTGRES_DSN is required")
	}
}
