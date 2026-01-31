package model

type PositionHistory struct {
	Id                int64   `json:"id"`
	CelestialObjectId int64   `json:"-"`
	Time              float64 `json:"time"`
	Position          Vector  `json:"position"`
}
