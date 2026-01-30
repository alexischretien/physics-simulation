package db

import (
	"fmt"
	"physics.simulation/model"
)

func (db *Database) CelestialObjectsBySimulationByID(id int64) ([]model.CelestialObject, error) {
	var celestialObjects []model.CelestialObject

	rows, err := db.Query("SELECT id, simulation_id, name, mass, x_position, y_position, z_position, x_velocity, y_velocity, z_velocity FROM celestial_object where simulation_id = ?", id)
	if err != nil {
		err = fmt.Errorf("celestialObjectsBySimulationId %d: %v", id, err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row model.CelestialObjectRow
		if err := rows.Scan(&row.Id, &row.SimulationId, &row.Name, &row.Mass, &row.X_position, &row.Y_position, &row.Z_position, &row.X_velocity, &row.Y_velocity, &row.Z_velocity); err != nil {
			err = fmt.Errorf("celestialObjectsBySimulationId %d: %v", id, err)
			return nil, err
		}
		celestialObjects = append(celestialObjects, model.NewCelestialObjectFromRow(row))
	}
	if err := rows.Err(); err != nil {
		err = fmt.Errorf("celestialObjectsBySimulationId %d: %v", id, err)
		return nil, err
	}
	return celestialObjects, nil
}
