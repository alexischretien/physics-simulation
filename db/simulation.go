package db

import (
	"fmt"

	"physics.simulation/model"
)

func (db *Database) SimulationExists(id int64) (bool, error) {
	var nbr int64
	row := db.QueryRow("SELECT COUNT(*) FROM simulation WHERE id = ?", id)
	if err := row.Scan(&nbr); err != nil {
		err = fmt.Errorf("SimulationExists: %v", err)
		return false, err
	}
	return nbr != 0, nil
}

func (db *Database) Simulations() ([]model.Simulation, error) {
	var sims []model.Simulation
	var sim model.Simulation
	rows, err := db.Query("SELECT id, title, duration, delta_t, writing_rate, is_dirty FROM simulation")
	if err != nil {
		err = fmt.Errorf("Simulations: %v", err)
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&sim.Id, &sim.Title, &sim.Duration, &sim.Delta_t, &sim.WritingRate, &sim.IsDirty); err != nil {
			err = fmt.Errorf("Simulations: %v", err)
			return nil, err
		}
		sims = append(sims, sim)
	}
	return sims, nil
}

func (db *Database) SimulationById(id int64) (*model.Simulation, error) {
	var sim model.Simulation
	row := db.QueryRow("SELECT id, title, duration, delta_t, writing_rate, is_dirty FROM simulation WHERE id = ?", id)
	if err := row.Scan(&sim.Id, &sim.Title, &sim.Duration, &sim.Delta_t, &sim.WritingRate, &sim.IsDirty); err != nil {
		err = fmt.Errorf("simulationByID %d: %v", id, err)
		return nil, err
	}
	return &sim, nil
}

func (db *Database) SimulationByIdLeftJoinChildrenTables(id int64) (*model.Simulation, error) {
	query := `SELECT s.id, s.title, s.duration, s.delta_t, s.writing_rate, s.is_dirty, 
	o.id, o.name, o.mass, o.x_position, o.y_position, o.z_position, o.x_velocity, o.y_velocity, o.z_velocity,
	p.id, p.time, p.x, p.y, p.z
	FROM simulation s
	LEFT JOIN celestial_object o on o.simulation_id = s.id
	LEFT JOIN position_history p on p.celestial_object_id = o.id
	WHERE s.id = ?`

	rows, err := db.Query(query, id)
	if err != nil {
		err = fmt.Errorf("SimulationByIdLeftJoinChildrenTables %d: %v", id, err)
		return nil, err
	}

	var sim model.SimulationRow
	var simulation model.Simulation
	celestialObjects := make(map[int64]model.CelestialObject)
	positionHistories := make(map[int64][]model.PositionHistory)
	firstRow := true
	for rows.Next() {
		err := rows.Scan(&sim.Id, &sim.Title, &sim.Duration, &sim.Delta_t, &sim.WritingRate, &sim.IsDirty,
			&sim.CelestialObjectId, &sim.Name, &sim.Mass, &sim.X_position, &sim.Y_position, &sim.Z_position, &sim.X_velocity, &sim.Y_velocity, &sim.Z_velocity,
			&sim.PositionHistoryId, &sim.Time, &sim.X, &sim.Y, &sim.Z)
		if err != nil {
			err = fmt.Errorf("SimulationByIdLeftJoinChildrenTables %d: %v", id, err)
			return nil, err
		}

		if firstRow {
			simulation = model.Simulation{
				Id: sim.Id, Title: sim.Title, Duration: sim.Duration, Delta_t: sim.Delta_t,
				WritingRate: sim.WritingRate, IsDirty: sim.IsDirty, CelestialObjects: []model.CelestialObject{},
			}
		}
		if sim.CelestialObjectId != nil {
			celestialObject, ok := celestialObjects[*sim.CelestialObjectId]
			if !ok {
				celestialObject = model.CelestialObject{
					Id: *sim.CelestialObjectId, SimulationId: sim.Id, Name: *sim.Name, Mass: *sim.Mass,
					Position:        model.Vector{X: *sim.X_position, Y: *sim.Y_position, Z: *sim.Z_position},
					Velocity:        model.Vector{X: *sim.X_velocity, Y: *sim.Y_velocity, Z: *sim.Z_velocity},
					PositionHistory: []model.PositionHistory{},
				}

				celestialObjects[*sim.CelestialObjectId] = celestialObject
			}
			if sim.PositionHistoryId != nil {
				positionHistory := model.PositionHistory{Id: *sim.PositionHistoryId, CelestialObjectId: *sim.CelestialObjectId,
					Time: *sim.Time, Position: model.Vector{X: *sim.X, Y: *sim.Y, Z: *sim.Z},
				}
				positionHistories[*sim.CelestialObjectId] = append(positionHistories[*sim.CelestialObjectId], positionHistory)
			}
		}
		firstRow = false
	}
	for celestialObjectId, celestialObject := range celestialObjects {
		celestialObject.PositionHistory = positionHistories[celestialObjectId]
		simulation.CelestialObjects = append(simulation.CelestialObjects, celestialObject)
	}
	return &simulation, nil
}

func (db *Database) UpdateSimulationIsDirty(id int64, isDirty bool) error {

	_, err := db.Query("UPDATE simulation set is_dirty = ? WHERE id = ?", isDirty, id)
	if err != nil {
		err = fmt.Errorf("UpdateSimulationIsDirty %d: %v", id, err)
		return err
	}
	return nil
}
