package selection

import (
	"math/rand"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
	"github.com/stretchr/testify/assert"
)

func TestTournamentSelectionNaturalFitness(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	p, _ := number.NewProbability(0.7)
	selector, err := NewTournamentSelection(
		WithConstantSelectionProbability(p))

	if assert.NoError(t, err) {
		steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)
		mary, _ := framework.NewEvaluatedCandidate("Mary", 9.1)
		john, _ := framework.NewEvaluatedCandidate("John", 8.4)
		gary, _ := framework.NewEvaluatedCandidate("Gary", 6.2)
		population := framework.EvaluatedPopulation{steve, mary, john, gary}

		// Run several iterations so that we get different tournament outcomes.
		for i := 0; i < 20; i++ {
			selection := selector.Select(population, true, 2, rng)
			assert.Len(t, selection, 2, "want selection size = 2, got", len(selection))
		}
	}
}

func TestTournamentSelectionNonNaturalFitness(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	p, _ := number.NewProbability(0.7)
	selector, err := NewTournamentSelection(
		WithConstantSelectionProbability(p))

	if assert.NoError(t, err) {
		gary, _ := framework.NewEvaluatedCandidate("Gary", 6.2)
		john, _ := framework.NewEvaluatedCandidate("John", 8.4)
		mary, _ := framework.NewEvaluatedCandidate("Mary", 9.1)
		steve, _ := framework.NewEvaluatedCandidate("Steve", 10.0)
		population := framework.EvaluatedPopulation{gary, john, mary, steve}

		// Run several iterations so that we get different tournament outcomes.
		for i := 0; i < 20; i++ {
			selection := selector.Select(population, false, 2, rng)
			assert.Len(t, selection, 2, "want selection size = 2, got", len(selection))
		}
	}
}

// The probability of selecting the fitter of two candidates must be greater
// than 0.5 to be useful (if it is not, there is no selection pressure, or the
// pressure is in favour of weaker candidates, which is counter-productive).
// This test ensures that an appropriate exception is thrown if the probability
// is 0.5 or less. Not throwing an exception is an error because it permits
// undetected bugs in evolutionary programs.
func TestTournamentSelectionProbabilityTooLow(t *testing.T) {
	ts, err := NewTournamentSelection(
		WithConstantSelectionProbability(number.ProbabilityEven))
	if assert.Error(t, err) {
		assert.Nil(t, ts)
	}
}