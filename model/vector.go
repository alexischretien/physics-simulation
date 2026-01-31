package model

import "math"

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func NewVector(x float64, y float64, z float64) Vector {
	return Vector{x, y, z}
}

func (v Vector) Copy() Vector {
	return Vector{v.X, v.Y, v.Z}
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector) UnitVector() Vector {
	return v.Divide(v.Magnitude())
}

func (v1 Vector) Add(v2 Vector) Vector {
	return Vector{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z}
}

func (v1 Vector) Substract(v2 Vector) Vector {
	return Vector{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
}

func (v1 Vector) Multiply(n float64) Vector {
	return Vector{v1.X * n, v1.Y * n, v1.Z * n}
}

func (v1 Vector) Divide(n float64) Vector {
	return Vector{v1.X / n, v1.Y / n, v1.Z / n}
}
