package steer

// Character is an autonomous character
type Character struct {
	Vehicle    Vehicle
	Behavior   Behavior
	Constraint Constraint
}

// Tick applies the steering behavior and advances the vehicle
func (c *Character) Tick() {
	c.Vehicle.Tick(c.Behavior())
	c.Vehicle.Position = c.Constraint(c.Vehicle.Position)
}
