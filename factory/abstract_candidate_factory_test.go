package factory

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/base"
	"github.com/stretchr/testify/assert"
)

type intFactory struct {
	AbstractCandidateFactory
}

func newIntFactory() *intFactory {
	return &intFactory{
		AbstractCandidateFactory{
			intGenerator{}},
	}
}

type intGenerator struct{}

func (ig intGenerator) GenerateRandomCandidate(rng *rand.Rand) base.Candidate { return rng.Int() }

func TestAbstractCandidateFactoryPopulationCreation(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	t.Run("generate whole population", func(t *testing.T) {
		cf := newIntFactory()
		pop := cf.GenerateInitialPopulation(10, rng)
		assert.Len(t, pop, 10)
	})

	t.Run("seed initial population", func(t *testing.T) {
		cf := newIntFactory()

		// preseed 5 candidates over 10
		preseed := make([]base.Candidate, 5)
		for i := range preseed {
			preseed[i] = i
		}

		pop := cf.SeedInitialPopulation(10, preseed, rng)
		assert.Len(t, pop, 10)
	})

	t.Run("too many seed candidates", func(t *testing.T) {
		cf := newIntFactory()

		// preseed 10 candidates
		preseed := make([]base.Candidate, 10)
		for i := range preseed {
			preseed[i] = i
		}

		assert.Panics(t, func() {
			cf.SeedInitialPopulation(5, preseed, rng)
		})
	})
}