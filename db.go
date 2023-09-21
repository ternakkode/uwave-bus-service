package main

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"

	_ "github.com/lib/pq"
)

var dbInstance *bun.DB

func InitDB(databaseDSN string) error {
	sqlDb, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		return err
	}

	dbInstance = bun.NewDB(sqlDb, pgdialect.New())
	dbInstance.AddQueryHook(
		bundebug.NewQueryHook(
			bundebug.WithEnabled(true),
			bundebug.WithVerbose(true),
		),
	)

	return nil
}

func GetDB() *bun.DB {
	return dbInstance
}
