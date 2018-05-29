package steer

import (
	"math"
	"math/rand"

	"github.com/jonahs99/vec"
)

// Behavior is a source of steering
type Behavior func() vec.Vec

// WanderBehavior returns a random-wander type behavior
func WanderBehavior(v *Vehicle, forward float64, radius float64, displacement float64) Behavior {
	dir := rand.Float64() * 2 * math.Pi
	return func() vec.Vec {
		dir += (2*rand.Float64() - 1) * displacement
		steer := vec.NewPolar(radius, dir)
		steer.Add(vec.NewPolar(forward, v.Heading))
		return steer
	}
}

// SeekBehavior returns a seek type behavior
func SeekBehavior(v *Vehicle, target *vec.Vec) Behavior {
	return func() vec.Vec {
		desiredVelocity := vec.Sub(*target, v.Position)
		desiredVelocity.Unit().Times(v.MaxSpeed)
		steer := vec.Sub(desiredVelocity, v.Velocity)
		return steer
	}
}

// FleeBehavior returns a flee type behavior
func FleeBehavior(v *Vehicle, target *vec.Vec) Behavior {
	return func() vec.Vec {
		desiredVelocity := vec.Sub(v.Position, *target)
		desiredVelocity.Unit().Times(v.MaxSpeed)
		steer := vec.Sub(desiredVelocity, v.Velocity)
		return steer
	}
}

// SeparationBehavior returns a separation type behavior
func SeparationBehavior(v *Vehicle, bunch *CharacterBunch, radius float64) Behavior {
	return func() vec.Vec {
		others := bunch.InRange(v, radius)
		steer := vec.Zero()
		for i := 0; i < len(others); i++ {
			other := &others[i].Vehicle
			radial := vec.Sub(v.Position, other.Position)
			dist2 := vec.Mag2(radial)
			if dist2 < 0.00001 {
				continue
			}
			weight := radius * radius / 16 / dist2
			radial.Unit().Times(weight)
			steer.Add(radial)
		}
		return steer
	}
}

// AlignBehavior returns an align type behavior
func AlignBehavior(v *Vehicle, bunch *CharacterBunch, radius float64) Behavior {
	return func() vec.Vec {
		others := bunch.InRange(v, radius)
		steer := vec.Zero()
		for i := 0; i < len(others); i++ {
			other := &others[i].Vehicle
			desiredVelocity := other.Velocity
			desiredVelocity.Unit().Times(v.MaxSpeed)
			steer.Add(vec.Sub(desiredVelocity, v.Velocity))
		}
		if len(others) > 1 {
			steer.Div(float64(len(others)))
		}
		return steer
	}
}

// CohesionBehavior returns a cohesion type behavior that move towards the center of mass
func CohesionBehavior(v *Vehicle, bunch *CharacterBunch, radius float64) Behavior {
	return func() vec.Vec {
		others := bunch.InRange(v, radius)
		if len(others) == 0 {
			return vec.Zero()
		}
		com := vec.Zero()
		for i := 0; i < len(others); i++ {
			other := &others[i].Vehicle
			com.Add(other.Position)
		}
		com.Div(float64(len(others)))
		steer := vec.Sub(com, v.Position)
		steer.Unit().Times(v.MaxSpeed)
		return steer
	}
}

// Linear returns a behavior that is a linear combination of other behaviors
func Linear(behaviors []Behavior, weights []float64) Behavior {
	if len(behaviors) != len(weights) {
		panic("len behaviors must equal len weights!")
	}
	return func() vec.Vec {
		steer := vec.Zero()
		for i := 0; i < len(behaviors); i++ {
			steerAtom := behaviors[i]()
			steerAtom.Times(weights[i])
			steer.Add(steerAtom)
		}
		return steer
	}
}

// Priority returns a behavior that is the first non-zero steering behavior
func Priority(behaviors ...Behavior) Behavior {
	epsilon := 0.001
	return func() vec.Vec {
		var steer vec.Vec
		for i := 0; i < len(behaviors); i++ {
			steer = behaviors[i]()
			if vec.Mag2(steer) > epsilon {
				break
			}
		}
		return steer
	}
}
