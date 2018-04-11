package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jonahs99/vec"

	"github.com/jonahs99/draws"
	"github.com/jonahs99/steering/steer"
)

func drawVehicle(context draws.Context, v *steer.Vehicle) {
	size := 16.0

	context.Save()
	context.Translate(v.Position.X, v.Position.Y)
	context.Rotate(v.Heading)

	context.BeginPath()
	context.MoveTo(size, 0)
	context.LineTo(-size/2, -size/2)
	context.LineTo(-size/2, size/2)
	context.ClosePath()

	context.FillStroke()

	context.Restore()
}

func wander(v *steer.Vehicle) func() vec.Vec {
	return func() vec.Vec {
		dir := rand.Float64() * 2 * math.Pi
		steer := vec.NewPolar(10, dir)
		steer.Add(vec.NewPolar(16, v.Heading))
		return steer
	}
}

func app(context draws.Context, quit <-chan struct{}) {
	context.Size(600, 600)
	context.BackgroundStyle("#eee")

	context.TranslateCenter()

	agents := make([]steer.Character, 0)

	for i := 0; i < 32; i++ {
		agents = append(agents, steer.Character{})
		a := &agents[len(agents)-1]

		a.Vehicle = steer.NewVehicle(rand.Float64()*300, rand.Float64()*300, 1, rand.Float64()*2*math.Pi)
		a.Behavior = wander(&a.Vehicle)
	}

	ticker := time.NewTicker(16 * time.Millisecond)

	context.StrokeStyle("#555")
	context.FillStyle("#06d")
	context.LineWidth(3)
	context.LineJoin("round")

	for {
		select {
		case <-ticker.C:
			context.Batch()
			context.Rect(-300, -300, 600, 600)
			context.Clear()

			for i := 0; i < len(agents); i++ {
				a := &agents[i]

				a.Tick()
				if a.Vehicle.Position.X < -300 {
					a.Vehicle.Position.X += 600
				} else if a.Vehicle.Position.X > 300 {
					a.Vehicle.Position.X -= 600
				}
				if a.Vehicle.Position.Y < -300 {
					a.Vehicle.Position.Y += 600
				} else if a.Vehicle.Position.Y > 300 {
					a.Vehicle.Position.Y -= 600
				}

				drawVehicle(context, &a.Vehicle)
			}

			context.Draw()
		case <-quit:
			fmt.Println("QUIT")
			return
		}
	}
}

func main() {
	draws.Serve(app, ":3000")
}
