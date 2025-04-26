package sql

import _ "embed"

var (
	//go:embed init.sql
	RecreateDatabaseSQL string
)
