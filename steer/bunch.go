package steer

import (
	"fmt"
	"math"

	"github.com/jonahs99/vec"
)

const gridSize = 24
const maxInGrid = 32

// CharacterBunch is a space-queryable group of characters
type CharacterBunch struct {
	Characters            []Character
	Width, Height         float64
	gridWidth, gridHeight int
	bins                  [gridSize * gridSize][maxInGrid]*Character
}

// NewCharacterBunch returns a new CharacterBunch
func NewCharacterBunch(characters []Character, w, h float64) CharacterBunch {
	fmt.Println(int(math.Ceil(w/float64(gridSize))), int(math.Ceil(h/float64(gridSize))))

	return CharacterBunch{
		characters,
		w, h,
		int(math.Ceil(w / float64(gridSize))), int(math.Ceil(h / float64(gridSize))),
		[gridSize * gridSize][maxInGrid]*Character{},
	}
}

// Tick ticks all the characters and updates the bins
func (cb *CharacterBunch) Tick() {
	for i := 0; i < gridSize*gridSize; i++ {
		cb.bins[i][0] = nil
	}

	for i := 0; i < len(cb.Characters); i++ {
		c := &cb.Characters[i]
		col, row := cb.gridCoords(c.Vehicle.Position)
		bin := &cb.bins[gridIndex(col, row)]

		for j := 0; j < maxInGrid-1; j++ {
			if bin[j] == nil {
				bin[j] = c
				bin[j+1] = nil
				break
			}
		}
	}

	for i := 0; i < len(cb.Characters); i++ {
		cb.Characters[i].Tick()
	}
}

// InRange returns the characters in the given radius
func (cb *CharacterBunch) InRange(v *Vehicle, radius float64) []*Character {
	characters := make([]*Character, 0)
	center := v.Position
	radius2 := radius * radius

	diag := vec.NewXY(radius, radius)
	min := vec.Sub(center, diag)
	max := vec.Add(center, diag)

	minCol, minRow := cb.gridCoords(min)
	maxCol, maxRow := cb.gridCoords(max)

	for col := minCol; col <= maxCol; col++ {
		for row := minRow; row <= maxRow; row++ {
			//fmt.Println(col, row, gridIndex(col, row))
			bin := &cb.bins[gridIndex(col, row)]
			for i := 0; i < maxInGrid; i++ {
				c := bin[i]
				if c == nil {
					break
				}
				if v == &c.Vehicle {
					continue
				}
				dist2 := vec.Mag2(vec.Sub(center, c.Vehicle.Position))
				if dist2 <= radius2 {
					characters = append(characters, c)
				}
			}
		}
	}

	return characters
}

func (cb *CharacterBunch) gridCoords(v vec.Vec) (int, int) {
	col := int(v.X) / cb.gridWidth
	row := int(v.Y) / cb.gridHeight
	col = clampInt(col, 0, gridSize-1)
	row = clampInt(row, 0, gridSize-1)
	return col, row
}

func gridIndex(col, row int) int {
	return col + row*gridSize
}

func clampInt(a, min, max int) int {
	if a < min {
		return min
	}
	if a > max {
		return max
	}
	return a
}
