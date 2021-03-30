package mutation

import (
	"math/rand"
	"testing"

	"github.com/arl/evolve/generator"
	"github.com/stretchr/testify/assert"
)

func TestListOrderMutation(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	cand := []int{1, 2, 3, 4, 5}

	population := []interface{}{cand}

	op := ListOrder{
		Count:  generator.ConstInt(1),
		Amount: generator.ConstInt(1),
	}
	mutpop := op.Apply(population, rng)

	assert.Len(t, mutpop, len(population), "population size should be unchanged after mutation")

	mutant := mutpop[0].([]int)
	t.Logf("original: %+v", cand)
	t.Logf("mutant : %+v", mutant)

	assert.Len(t, mutant, len(cand), "mutant should be same length as original")

	// The mutant should have the same elements but in a different order
	matches := 0
	for i := range cand {
		if cand[i] == mutant[i] {
			matches++
		} else {
			// If positions don't match, an adjacent item should be a match.
			next := (i + 1) % len(cand)
			prev := ((i - 1) + len(cand)) % len(cand)
			matchAdjacent := (cand[i] == mutant[next]) != (cand[i] == mutant[prev])
			assert.True(t, matchAdjacent, "Mutated item is not in one of the expected positions")
		}
	}

	assert.Equal(t, matches, len(cand)-2, "All but 2 positions should be unchanged.")
}
