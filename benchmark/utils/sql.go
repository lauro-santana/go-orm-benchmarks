package utils

import _ "embed"

var (
	//go:embed sql/insert.sql
	InsertQuery string
	//go:embed sql/insert_returning_id.sql
	InsertReturningIDQuery string
	//go:embed sql/insert_bulk.sql
	InsertBulkQuery string
	//go:embed sql/update.sql
	UpdateQuery string
	//go:embed sql/delete.sql
	DeleteQuery string
	//go:embed sql/select_by_id.sql
	SelectByIDQuery string
	//go:embed sql/select_paginating.sql
	SelectPaginatingQuery string
)
