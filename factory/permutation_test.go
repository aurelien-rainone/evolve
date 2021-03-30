package factory

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListPermutation(t *testing.T) {
	elements := make([]int, 10)

	for i := 0; i < len(elements); i++ {
		elements[i] = i
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	f := ListPermutation(elements)
	cpy := f.New(rng)

	assert.Len(t, cpy, len(elements), "permutated list should have the same length")

	same := reflect.ValueOf(elements).Pointer() == reflect.ValueOf(cpy).Pointer()
	assert.False(t, same, "new individual has the same backing array as the original")

	assert.NotEqualValues(t, elements, cpy, "new indivual should have permuted values")
}
