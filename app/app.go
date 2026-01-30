package app

import (
	"physics.simulation/db"
)

type App struct {
	Database *db.Database
}

func New() (app *App, err error) {
	dbConfig, err := db.InitConfig()
	if err != nil {
		return nil, err
	}
	database, err := db.New(dbConfig)
	if err != nil {
		return nil, err
	}
	return &App{Database: database}, err
}

func (a *App) Close() error {
	return a.Database.Close()
}
