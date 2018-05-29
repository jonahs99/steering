package steer

import (
	"github.com/jonahs99/vec"
)

// NewVehicle returns a vehicle with some default parameters
func NewVehicle(x, y, speed, heading float64) Vehicle {
	return Vehicle{
		Mass:     1,
		Position: vec.NewXY(x, y),
		Velocity: vec.NewPolar(speed, heading),
		Heading:  heading,
		MaxForce: 0.2,
		MaxSpeed: 3,
	}
}

// Vehicle is a point Mass physical approximation
type Vehicle struct {
	Mass     float64
	Position vec.Vec
	Velocity vec.Vec
	Heading  float64

	MaxForce float64
	MaxSpeed float64
}

// Tick is one simulation step
func (v *Vehicle) Tick(steeringDirection vec.Vec) {
	steeringForce := vec.Truncate(steeringDirection, v.MaxForce)
	acceleration := vec.Times(steeringForce, 1/v.Mass)
	v.Velocity.Add(acceleration)
	v.Velocity = vec.Truncate(v.Velocity, v.MaxSpeed)
	v.Position.Add(v.Velocity)
	v.Heading = vec.Theta(v.Velocity)
}
