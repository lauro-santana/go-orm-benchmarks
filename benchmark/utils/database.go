package utils

import (
	"database/sql"
	"log"

	queries "github.com/lauro-santana/golang-orm-benchmarks/sql"

	// Postgres driver.
	_ "github.com/jackc/pgx/v5/stdlib"
)

func RecreateDatabase() {
	db, err := sql.Open("pgx", PostgresDSN)
	if err != nil {
		log.Fatal("the benchmark execution was aborted", err)
	}

	defer func() {
		_ = db.Close()
	}()

	_, err = db.Exec(queries.RecreateDatabaseSQL)
	if err != nil {
		log.Fatal("the benchmark execution was aborted", err)
	}
}
