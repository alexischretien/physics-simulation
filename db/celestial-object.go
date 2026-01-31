package db

import (
	"physics.simulation/model"
)

func (db *Database) CelestialObjectsBySimulationByID(id int64) ([]model.CelestialObject, error) {
	var celestialObjects []model.CelestialObject

	rows, err := db.Query("SELECT id, simulation_id, name, mass, x_position, y_position, z_position, x_velocity, y_velocity, z_velocity FROM celestial_object where simulation_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row model.CelestialObjectRow
		if err := rows.Scan(&row.Id, &row.SimulationId, &row.Name, &row.Mass, &row.X_position, &row.Y_position, &row.Z_position, &row.X_velocity, &row.Y_velocity, &row.Z_velocity); err != nil {
			return nil, err
		}
		celestialObjects = append(celestialObjects, model.NewCelestialObjectFromRow(row))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return celestialObjects, nil
}

func (db *Database) CreateCelestialObjectForSimulation(simulationId int64, c model.CelestialObject) (*int64, error) {
	res, err := db.Exec(
		`INSERT INTO celestial_object (simulation_id, name, mass, x_position, y_position, z_position, x_velocity, y_velocity, z_velocity)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		simulationId, c.Name, c.Mass, c.Position.X, c.Position.Y, c.Position.Z, c.Velocity.X, c.Velocity.Y, c.Velocity.Z)

	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (db *Database) UpdateCelestialObject(c model.CelestialObject) error {
	_, err := db.Exec(
		`UPDATE celestial_object 
		SET name = ?, mass = ?, x_position = ?, y_position = ?, z_position = ?, x_velocity = ?, y_velocity = ?, z_velocity = ?
		WHERE id = ?`,
		c.Name, c.Mass, c.Position.X, c.Position.Y, c.Position.Z, c.Velocity.X, c.Velocity.Y, c.Velocity.Z, c.Id)

	if err != nil {
		return err
	}
	return nil
}

func (db *Database) DeleteCelestialObject(id int64) error {
	_, err := db.Query("DELETE FROM celestial_object WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}
