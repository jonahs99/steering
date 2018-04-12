package steer

import "github.com/jonahs99/vec"

// CharacterBunch is a space-queryable group of characters
type CharacterBunch struct {
	Characters []Character
}

// InRange returns the characters in the given radius
func (cb *CharacterBunch) InRange(center vec.Vec, radius float64) []*Character {
	characters := make([]*Character, 0)
	radius2 := radius * radius
	for i := 0; i < len(cb.Characters); i++ {
		c := &cb.Characters[i]
		dist2 := vec.Mag2(vec.Sub(center, c.Vehicle.Position))
		if dist2 <= radius2 {
			characters = append(characters, c)
		}
	}
	return characters
}
