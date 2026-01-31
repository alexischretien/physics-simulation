package db

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	*sql.DB
}

func New(config *Config) (*Database, error) {
	cfg := mysql.NewConfig()
	cfg.User = config.DbUser
	cfg.Passwd = config.DbPassword
	cfg.Addr = config.DbAddress
	cfg.DBName = config.DbName
	cfg.Net = "tcp"

	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		return nil, fmt.Errorf("Db.new: could not connect to the database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Db.new: failed to ping the database: %v", err)
	}
	return &Database{db}, nil
}
