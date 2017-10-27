package operators

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/bitstring"
	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
	"github.com/stretchr/testify/assert"
)

// Ensures that mutation occurs correctly. Because of the random aspect we
// can't actually make many assertions. This just ensures that there are no
// unexpected exceptions and that the length of the bit strings remains as
// expected.
func TestBitStringMutationRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	var (
		mutation *AbstractMutation
		original *bitstring.BitString
		err      error
	)
	mutation, err = NewBitStringMutation(ConstantProbability(number.ProbabilityEven))
	assert.NoError(t, err)

	original, err = bitstring.NewFromString("111100101")
	assert.NoError(t, err)

	population := []framework.Candidate{original}
	for i := 0; i < 20; i++ {
		// Perform several iterations to get different mutations.
		population = mutation.Apply(population, rng)
		mutated := population[0]
		assert.IsType(t, &bitstring.BitString{}, mutated)
		ms := mutated.(*bitstring.BitString)
		assert.Equalf(t, 9, ms.Len(), "want mutated bit string to not change length, 9, got %v", ms.Len())
	}
}

// Ensures that mutation occurs correctly.  Uses a probability of 1 to
// make the outcome predictable (all bits will be flipped).
func TestBitStringMutationSingleBit(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	var (
		mutation *AbstractMutation
		original *bitstring.BitString
		err      error
	)

	mutation, err = NewBitStringMutation()
	assert.NoError(t, err)

	original, err = bitstring.NewFromString("111100101")
	assert.NoError(t, err)

	population := []framework.Candidate{original}
	population = mutation.Apply(population, rng)

	mutated := population[0]
	assert.IsType(t, &bitstring.BitString{}, mutated)
	ms := mutated.(*bitstring.BitString)

	assert.False(t, ms.Equals(original), "want mutant to be different from original, got equals")
	assert.Equalf(t, 9, ms.Len(), "want mutated bit string to not change length, 9, got %v", ms.Len())
	set := ms.CountSetBits()
	unset := ms.CountUnsetBits()
	assert.Truef(t, set == 5 || set == 7, "want 5 or 7 set bits in mutated bit string, got %v", set)
	assert.Truef(t, unset == 2 || unset == 4, "want 2 or 4 unset bits in mutated bit string, got %v", unset)
}

// TODO:  test BitStringMutation constructed with options other than default ones
//func