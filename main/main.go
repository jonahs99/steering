package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jonahs99/draws"
	"github.com/jonahs99/steering/steer"
)

func initContext(context draws.Context, w, h float64) {
	context.Size(w, h)
	context.BackgroundStyle("#eee")

	context.StrokeStyle("#555")
	context.FillStyle("#06d")
	context.LineWidth(3)
	context.LineJoin("round")
}

func drawVehicle(context draws.Context, v *steer.Vehicle) {
	size := 9.0

	context.Save()
	context.Translate(v.Position.X, v.Position.Y)
	context.Rotate(v.Heading)

	context.BeginPath()
	context.MoveTo(size, 0)
	context.LineTo(-size/2, -size/2)
	context.LineTo(-size/2, size/2)
	context.ClosePath()

	context.Fill()

	context.Restore()

	/*context.BeginPath()
	context.Circle(v.Position.X, v.Position.Y, size/2)
	context.Fill()*/
}

func app(context draws.Context, quit <-chan struct{}) {
	w, h := 1000.0, 1000.0

	n := 128
	bunch := steer.NewCharacterBunch(make([]steer.Character, n), w, h)

	for i := 0; i < n; i++ {
		a := &bunch.Characters[i]
		a.Vehicle = steer.NewVehicle(rand.Float64()*w, rand.Float64()*h, 1, rand.Float64()*2*math.Pi)
		a.Constraint = steer.WrapConstraint(0, 0, w, h)

		flockRadius := 64.0

		separation := steer.SeparationBehavior(&a.Vehicle, &bunch, flockRadius)
		align := steer.AlignBehavior(&a.Vehicle, &bunch, flockRadius)
		cohesion := steer.CohesionBehavior(&a.Vehicle, &bunch, flockRadius)
		flock := steer.Linear([]steer.Behavior{separation, align, cohesion}, []float64{0.7, 0.25, 0.1})
		wander := steer.WanderBehavior(&a.Vehicle, 12, 8, 0.1)

		a.Behavior = steer.Linear([]steer.Behavior{flock, wander}, []float64{0.98, 0.02})
		//a.Behavior = steer.Priority(flock, wander)
	}

	initContext(context, w, h)

	frame := 0

	for ticker := time.NewTicker(16 * time.Millisecond); ; {
		select {
		case <-ticker.C:
			//start := time.Now()
			bunch.Tick()
			//tickTime := time.Since(start)

			if frame%1 == 0 {
				//start = time.Now()
				context.Batch()
				context.Clear()
				for i := 0; i < len(bunch.Characters); i++ {
					a := &bunch.Characters[i]
					drawVehicle(context, &a.Vehicle)
				}
				context.Draw()
				//drawTime := time.Since(start)

				//fmt.Printf("%v to tick, %v to draw\n", tickTime, drawTime)
			}
			frame++

		case <-quit:
			fmt.Println("QUIT")
			return
		}
	}
}

func main() {
	draws.Serve(app, ":3000")
}
