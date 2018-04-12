package steer

import (
	"math"
	"math/rand"

	"github.com/jonahs99/vec"
)

// Behavior is a source of steering
type Behavior func() vec.Vec

// WanderBehavior returns a random-wander type behavior
func WanderBehavior(v *Vehicle, forward float64, radius float64) Behavior {
	return func() vec.Vec {
		dir := rand.Float64() * 2 * math.Pi
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
