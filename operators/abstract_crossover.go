package operators

import (
	"fmt"
	"math/rand"

	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
)

// Mater is the interface implemented by objects defining the Mate function.
type Mater interface {

	// Mate performs crossover on a pair of parents to generate a pair of
	// offspring.
	//
	// parent1 and parent2 are the two individuals that provides the source
	// material for generating offspring.
	Mate(parent1, parent2 framework.Candidate,
		numberOfCrossoverPoints int64,
		rng *rand.Rand) []framework.Candidate
}

// AbstractCrossover is a generic struct for crossover implementations.
//
// It supports all crossover processes that operate on a pair of parent
// candidates.
// Both the number of crossovers points and the crossover probability are
// configurable. Crossover is applied to a proportion of selected parent pairs,
// with the remainder copied unchanged into the output population. The size of
// this evolved proportion is controlled by the code crossoverProbability
// parameter.
type AbstractCrossover struct {
	crossoverPointsVariable      number.IntegerGenerator
	crossoverProbabilityVariable number.ProbabilityGenerator
	Mater
}

// NewAbstractCrossover creates an AbstractCrossover configured with the
// provided options.
//
// TODO: example of use of how setting options
func NewAbstractCrossover(mater Mater, options ...Option) (*AbstractCrossover, error) {
	// create with default options, 1 crossover point with a probability of 1
	op := &AbstractCrossover{
		crossoverPointsVariable:      number.NewConstantIntegerGenerator(1),
		crossoverProbabilityVariable: number.NewConstantProbabilityGenerator(number.ProbabilityOne),
		Mater: mater,
	}

	// set client options
	for _, option := range options {
		if err := option.Apply(op); err != nil {
			return nil, fmt.Errorf("can't apply abstract crossover options: %v", err)
		}
	}
	return op, nil
}

// Apply applies the crossover operation to the selected candidates.
//
// Pairs of candidates are chosen randomly and subjected to crossover to
// produce a pair of offspring candidates.
//
// The selectedCandidates are the evolved individuals that have survived to be
// eligible to reproduce.
//
// It returns the combined set of evolved offspring generated by applying
// crossover to the the selected candidates.
func (op *AbstractCrossover) Apply(selectedCandidates []framework.Candidate, rng *rand.Rand) []framework.Candidate {
	// Shuffle the collection before applying each operation so that the
	// evolution is not influenced by any ordering artifacts from previous
	// operations.
	selectionClone := make([]framework.Candidate, len(selectedCandidates))
	copy(selectionClone, selectedCandidates)
	framework.ShuffleCandidates(selectionClone, rng)

	result := make([]framework.Candidate, 0, len(selectedCandidates))
	var iterator = 0
	for iterator < len(selectionClone) {
		parent1 := selectionClone[iterator]
		iterator++
		if iterator < len(selectionClone) {
			parent2 := selectionClone[iterator]
			iterator++
			// Randomly decide (according to the current crossover probability)
			// whether to perform crossover for these 2 parents.
			var crossoverPoints int64
			if op.crossoverProbabilityVariable.NextValue().NextEvent(rng) {
				crossoverPoints = op.crossoverPointsVariable.NextValue()
			}

			if crossoverPoints > 0 {
				result = append(result, op.Mate(parent1, parent2, crossoverPoints, rng)...)
			} else {
				// If there is no crossover to perform, just add the parents to the
				// results unaltered.
				result = append(result, parent1, parent2)
			}
		} else {
			// If we have an odd number of selected candidates, we can't pair up
			// the last one so just leave it unmodified.
			result = append(result, parent1)
		}
	}
	return result
}
