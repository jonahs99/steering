package steer

import "github.com/jonahs99/vec"

// Character is an autonomous character
type Character struct {
	Vehicle  Vehicle
	Behavior func() vec.Vec
}

// Tick applies the steering behavior and advances the vehicle
func (c *Character) Tick() {
	c.Vehicle.Tick(c.Behavior())
}
