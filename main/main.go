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
	context.TranslateCenter()

	context.StrokeStyle("#555")
	context.FillStyle("#06d")
	context.LineWidth(3)
	context.LineJoin("round")
}

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

func app(context draws.Context, quit <-chan struct{}) {
	w, h := 1200.0, 800.0

	n := 32
	agents := make([]steer.Character, n)
	for i := 0; i < n; i++ {
		a := &agents[i]

		a.Vehicle = steer.NewVehicle(rand.Float64()*w-w/2, rand.Float64()*h-h/2, 1, rand.Float64()*2*math.Pi)
		/*if i == 0 {
			a.Behavior = steer.SeekBehavior(&a.Vehicle, &agents[1].Vehicle.Position)
			a.Vehicle.Mass *= 2
			a.Constraint = steer.WrapConstraint(-w/2, -h/2, w/2, h/2)
			//a.Behavior = steer.SeekBehavior(&a.Vehicle, &agents[n-1].Vehicle.Position)
		} else if i < n-1 {*/
		a.Behavior = steer.Linear(
			[]steer.Behavior{steer.FleeBehavior(&a.Vehicle, &agents[(i-1+n)%n].Vehicle.Position),
				steer.SeekBehavior(&a.Vehicle, &agents[(i+1)%n].Vehicle.Position)},
			[]float64{0.5, 0.4},
		)
		a.Constraint = steer.NilConstraint()
		//steer.CircleWrapConstraint(vec.Zero(), h/2)
		/*} else {
			a.Behavior = steer.FleeBehavior(&a.Vehicle, &agents[i-1].Vehicle.Position)
			a.Constraint = steer.CircleWrapConstraint(vec.Zero(), h/2)
		}*/
	}

	initContext(context, w, h)

	for ticker := time.NewTicker(16 * time.Millisecond); ; {
		select {
		case <-ticker.C:
			context.Batch()
			context.Clear()

			for i := 0; i < len(agents); i++ {
				a := &agents[i]

				a.Tick()

				if i == 0 {
					context.FillStyle("#f06")
				} else if i == 1 {
					context.FillStyle("#0f6")
				} else if i == 2 {
					context.FillStyle("#60f")
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
