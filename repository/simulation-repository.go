package repository

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	. "physics.simulation/model"
	"strings"
)

var db *sql.DB

func ConnectToDb() {
	// Capture connection properties.
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "physics_simulation"

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

func Simulations() ([]Simulation, error, int) {
	var sims []Simulation
	var sim Simulation
	rows, err := db.Query("SELECT id, title, duration, delta_t, writing_rate, is_dirty FROM simulation")
	if err != nil {
		err = fmt.Errorf("Simulations: %v", err)
		return nil, err, http.StatusInternalServerError
	}
	for rows.Next() {
		err := rows.Scan(&sim.Id, &sim.Title, &sim.Duration, &sim.Delta_t, &sim.WritingRate, &sim.IsDirty)
		if err != nil {
			err = fmt.Errorf("Simulations: %v", err)
			return nil, err, http.StatusInternalServerError
		}
		sims = append(sims, sim)
	}
	return sims, nil, http.StatusOK
}

func SimulationByIdLeftJoinChildrenTables(id int64) (*Simulation, error, int) {
	_, err, status := SimulationById(id)
	if err != nil {
		return nil, err, status
	}

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
		return nil, err, http.StatusInternalServerError
	}

	var sim SimulationRow
	var simulation Simulation
	celestialObjects := make(map[int64]CelestialObject)
	positionHistories := make(map[int64][]PositionHistory)
	firstRow := true
	for rows.Next() {
		err := rows.Scan(&sim.Id, &sim.Title, &sim.Duration, &sim.Delta_t, &sim.WritingRate, &sim.IsDirty,
			&sim.CelestialObjectId, &sim.Name, &sim.Mass, &sim.X_position, &sim.Y_position, &sim.Z_position, &sim.X_velocity, &sim.Y_velocity, &sim.Z_velocity,
			&sim.PositionHistoryId, &sim.Time, &sim.X, &sim.Y, &sim.Z)
		if err != nil {
			err = fmt.Errorf("SimulationByIdLeftJoinChildrenTables %d: %v", id, err)
			return nil, err, http.StatusInternalServerError
		}

		if firstRow {
			simulation = Simulation{Id: sim.Id, Title: sim.Title, Duration: sim.Duration, Delta_t: sim.Delta_t, WritingRate: sim.WritingRate, IsDirty: sim.IsDirty, CelestialObjects: []CelestialObject{}}
		}
		if sim.CelestialObjectId != nil {
			celestialObject, ok := celestialObjects[*sim.CelestialObjectId]
			if !ok {
				celestialObject = CelestialObject{Id: *sim.CelestialObjectId, SimulationId: sim.Id, Name: *sim.Name, Mass: *sim.Mass, Position: Vector{X: *sim.X_position, Y: *sim.Y_position, Z: *sim.Z_position}, Velocity: Vector{X: *sim.X_velocity, Y: *sim.Y_velocity, Z: *sim.Z_velocity}, PositionHistory: []PositionHistory{}}
				celestialObjects[*sim.CelestialObjectId] = celestialObject
			}
			if sim.PositionHistoryId != nil {
				positionHistory := PositionHistory{Id: *sim.PositionHistoryId, CelestialObjectId: *sim.CelestialObjectId, Time: *sim.Time, Position: Vector{X: *sim.X, Y: *sim.Y, Z: *sim.Z}}
				positionHistories[*sim.CelestialObjectId] = append(positionHistories[*sim.CelestialObjectId], positionHistory)
			}
		}
		firstRow = false
	}
	for celestialObjectId, celestialObject := range celestialObjects {
		celestialObject.PositionHistory = positionHistories[celestialObjectId]
		simulation.CelestialObjects = append(simulation.CelestialObjects, celestialObject)
	}
	return &simulation, nil, http.StatusOK
}

func SimulationById(id int64) (*Simulation, error, int) {
	var sim Simulation
	row := db.QueryRow("SELECT id, title, duration, delta_t, writing_rate, is_dirty FROM simulation WHERE id = ?", id)
	if err := row.Scan(&sim.Id, &sim.Title, &sim.Duration, &sim.Delta_t, &sim.WritingRate, &sim.IsDirty); err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("simulationByID %d: no such simulation", id)
			return nil, err, http.StatusNotFound
		}
		err = fmt.Errorf("simulationByID %d: %v", id, err)
		return nil, err, http.StatusInternalServerError
	}
	return &sim, nil, http.StatusOK
}

func CelestialObjectsBySimulationByID(id int64) ([]CelestialObject, error, int) {
	_, err, status := SimulationById(id)
	if err != nil {
		return nil, err, status
	}

	var celestialObjects []CelestialObject

	rows, err := db.Query("SELECT id, simulation_id, name, mass, x_position, y_position, z_position, x_velocity, y_velocity, z_velocity FROM celestial_object where simulation_id = ?", id)
	if err != nil {
		err = fmt.Errorf("celestialObjectsBySimulationId %d: %v", id, err)
		return nil, err, http.StatusInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		var row CelestialObjectRow
		if err := rows.Scan(&row.Id, &row.SimulationId, &row.Name, &row.Mass, &row.X_position, &row.Y_position, &row.Z_position, &row.X_velocity, &row.Y_velocity, &row.Z_velocity); err != nil {
			err = fmt.Errorf("celestialObjectsBySimulationId %d: %v", id, err)
			return nil, err, http.StatusInternalServerError
		}
		celestialObjects = append(celestialObjects, NewCelestialObjectFromRow(row))
	}
	if err := rows.Err(); err != nil {
		err = fmt.Errorf("celestialObjectsBySimulationId %d: %v", id, err)
		return nil, err, http.StatusInternalServerError
	}
	return celestialObjects, nil, http.StatusOK
}

func SaveHistoryPositions(celestialObjects []CelestialObject) error {

	query := "INSERT INTO position_history(celestial_object_id, time, x, y, z) VALUES "
	var values []interface{}

	for _, o := range celestialObjects {
		for _, p := range o.PositionHistory {
			query += "(?, ?, ?, ?, ?),"
			values = append(values, o.Id, p.Time, p.Position.X, p.Position.Y, p.Position.Z)
		}
	}
	query = strings.TrimSuffix(query, ",")

	if len(values) > 0 {
		_, err := db.Exec(query, values...)
		if err != nil {
			err = fmt.Errorf("SaveHistoryPositions: %v", err)
			return err
		}
	}
	return nil
}

func UpdateSimulationIsDirty(id int64, isDirty bool) error {
	_, err := db.Query("UPDATE simulation set is_dirty = ? WHERE id = ?", isDirty, id)
	if err != nil {
		err = fmt.Errorf("UpdateSimulationIsDirty %d: %v", id, err)
		return err
	}
	return nil
}
