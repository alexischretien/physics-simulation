package model

type PositionHistory struct {
	Id                int64
	CelestialObjectId int64
	Time              float64
	Position          Vector
}
