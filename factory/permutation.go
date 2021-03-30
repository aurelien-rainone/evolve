package factory

import (
	"math/rand"
)

type ListPermutation []int

func (f ListPermutation) New(rng *rand.Rand) interface{} {
	// shuffle modifies the slice in-place, and since we want new individuals to
	// be created we have to create a new slice.
	cpy := make([]int, len(f))
	copy(cpy, f)
	rng.Shuffle(len(cpy), func(i, j int) {
		cpy[i], cpy[j] = cpy[j], cpy[i]
	})

	return cpy
}
