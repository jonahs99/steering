package steer

import (
	"math"

	"github.com/jonahs99/vec"
)

// Constraint is a function which may modify the tentative position
type Constraint func(position vec.Vec) vec.Vec

// NilConstraint returns a constraint that doesn't constrain
func NilConstraint() Constraint {
	return func(position vec.Vec) vec.Vec { return position }
}

// RectBoundConstraint returns a constraint that keeps an agents in the world
func RectBoundConstraint(minx, miny, maxx, maxy float64) Constraint {
	return func(position vec.Vec) vec.Vec {
		return clampXY(position, minx, maxx, miny, maxy)
	}
}

// CircleBoundConstraint returns a constraint that keep agents in the circle
func CircleBoundConstraint(center vec.Vec, radius float64) Constraint {
	return func(position vec.Vec) vec.Vec {
		return vec.Add(vec.Truncate(vec.Sub(position, center), radius), center)
	}
}

// WrapConstraint returns a constraint that wraps the character around a rectangle
func WrapConstraint(minx, miny, maxx, maxy float64) Constraint {
	return func(position vec.Vec) vec.Vec {
		if position.X < minx {
			position.X += maxx - minx
		} else if position.X > maxx {
			position.X -= maxx - minx
		}
		if position.Y < miny {
			position.Y += maxy - miny
		} else if position.Y > maxy {
			position.Y -= maxy - miny
		}
		return position
	}
}

// CircleWrapConstraint returns a constraint that wraps the agent around a circle
func CircleWrapConstraint(center vec.Vec, radius float64) Constraint {
	radius2 := radius * radius
	return func(position vec.Vec) vec.Vec {
		radial := vec.Sub(center, position)
		if vec.Mag2(radial) > radius2 {
			radial.Unit().Times(radius * 2)
			position.Add(radial)
		}
		return position
	}
}

func clamp(a, min, max float64) float64 {
	return math.Min(max, math.Max(min, a))
}

func clampXY(v vec.Vec, minx, maxx, miny, maxy float64) vec.Vec {
	return vec.NewXY(clamp(v.X, minx, maxx), clamp(v.Y, miny, maxy))
}
