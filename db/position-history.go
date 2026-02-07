package db

import (
	"strings"

	"physics.simulation/model"
)

func (db *Database) SavePositionHistories(celestialObjects []model.CelestialObject) error {

	query := "INSERT INTO position_history(celestial_object_id, time, x, y, z) VALUES "
	var values []any

	for _, o := range celestialObjects {
		for _, p := range o.PositionHistory {
			query += "(?, ?, ?, ?, ?),"
			values = append(values, o.Id, p.Time, p.Position.X, p.Position.Y, p.Position.Z)
		}
	}
	query = strings.TrimSuffix(query, ",")

	if len(values) > 0 {
		if _, err := db.Exec(query, values...); err != nil {
			return err
		}
	}
	return nil
}

func (db *Database) DeletePositionHistoryForSimulationId(id int64) error {
	_, err := db.Query("delete from position_history where celestial_object_id in (select id from celestial_object where simulation_id = ?)", id)
	if err != nil {
		return err
	}
	return nil
}
