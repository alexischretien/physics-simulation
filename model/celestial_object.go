package model

import (
	"fmt"
	"math"
	"os"
)

type CelestialObject struct {
	Id              int64
	SimulationId    int64
	Name            string
	Mass            float64
	Position        Vector
	Velocity        Vector
	Acceleration    Vector
	PositionHistory []PositionHistory
	dataFile        *os.File
}

type CelestialObjectRow struct {
	Id           int64
	SimulationId int64
	Name         string
	Mass         float64
	X_position   float64
	Y_position   float64
	Z_position   float64
	X_velocity   float64
	Y_velocity   float64
	Z_velocity   float64
}

func NewCelestialObject(id int64, name string, mass float64, position Vector, velocity Vector, acceleration Vector) CelestialObject {
	return CelestialObject{id, 0, name, mass, position, velocity, acceleration, []PositionHistory{}, nil}
}

func NewCelestialObjectFromRow(row CelestialObjectRow) CelestialObject {
	return NewCelestialObject(row.Id, row.Name, row.Mass, NewVector(row.X_position, row.Y_position, row.Z_position), NewVector(row.X_velocity, row.Y_velocity, row.Z_velocity), NewVector(0.0, 0.0, 0.0))
}

func (c *CelestialObject) UpdateAccelerationVelocityPosition(celestialObjects []CelestialObject, i int, dt float64) {

	var accGrav Vector = Vector{0.0, 0.0, 0.0}
	var distance Vector
	var distanceMagnitude float64
	var distanceUnitVector Vector

	for j := range celestialObjects {
		if i != j {
			distance = celestialObjects[j].Position.Substract(c.Position)
			distanceMagnitude = distance.Magnitude()
			distanceUnitVector = distance.UnitVector()
			accGrav = accGrav.Add(distanceUnitVector.Multiply(celestialObjects[j].Mass).Divide(distanceMagnitude * distanceMagnitude))
		}
	}
	accGrav = accGrav.Multiply(6.6740831e-11) // grav constant

	c.Acceleration = accGrav
	c.Velocity = c.Velocity.Add(c.Acceleration.Multiply(dt))
	c.Position = c.Position.Add(c.Velocity.Multiply(dt))
}

func (c1 CelestialObject) Equals(c2 CelestialObject) bool {
	return (((c1.Mass == 0.0 && c2.Mass == 0.0) || (math.Abs(c1.Mass-c2.Mass)/((c1.Mass+c2.Mass)/2) < 0.0001)) && c1.Acceleration == c2.Acceleration && c1.Velocity == c2.Velocity && c1.Position == c2.Position)
}

func (c *CelestialObject) setDataFile(file *os.File) {
	c.dataFile = file
}

func (c CelestialObject) CloseDataFile() {
	c.dataFile.Close()
}

func (c CelestialObject) AppendCurrentPositionToDataFile() {
	data := []byte(fmt.Sprintf("%v %v %v\n", c.Position.X, c.Position.Y, c.Position.Z))
	_, err := c.dataFile.Write(data)
	if err != nil {
		panic(err)
	}
	c.dataFile.Sync()
}
