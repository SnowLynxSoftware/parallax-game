package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/snowlynxsoftware/parallax-game/server/util"
)

type AppDataSource struct {
	DB *sqlx.DB
}

func NewAppDataSource() *AppDataSource {
	return &AppDataSource{}
}

func (d *AppDataSource) Connect(dbConnectionString string) {
	db, err := sqlx.Connect("postgres", dbConnectionString)
	if err != nil {
		panic(err)
	}
	d.DB = db
	err = d.DB.Ping()
	if err != nil {
		panic(err)
	}
	util.LogInfo("Connected to database")
}
