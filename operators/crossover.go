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
	// TODO: should return 2 values of a slice of 2 values
	Mate(parent1, parent2 framework.Candidate,
		numberOfCrossoverPoints int64,
		rng *rand.Rand) []framework.Candidate
}

// Crossover is a generic struct for crossover implementations.
//
// It supports all crossover processes that operate on a pair of parent
// candidates.
// Both the number of crossovers points and the crossover probability are
// configurable. Crossover is applied to a proportion of selected parent pairs,
// with the remainder copied unchanged into the output population. The size of
// this evolved proportion is controlled by the code crossoverProbability
// parameter.
type Crossover struct {
	npts number.IntegerGenerator
	prob number.ProbabilityGenerator
	Mater
}

// NewCrossover creates a standard Crossover configured with the
// provided options.
//
// TODO: example of use of how setting options
func NewCrossover(mater Mater, options ...Option) (*Crossover, error) {
	// create with default options, 1 crossover point with a probability of 1
	op := &Crossover{
		npts:  number.NewConstantIntegerGenerator(1),
		prob:  number.NewConstantProbabilityGenerator(number.ProbabilityOne),
		Mater: mater,
	}

	// set client options
	for _, option := range options {
		if err := option.Apply(op); err != nil {
			return nil, fmt.Errorf("can't apply crossover option: %v", err)
		}
	}
	return op, nil
}

// CrossoverPoints returns the crossover points generator.
func (op *Crossover) CrossoverPoints() number.IntegerGenerator {
	return op.npts
}

// CrossoverProbability returns the crossover probability generator.
func (op *Crossover) CrossoverProbability() number.ProbabilityGenerator {
	return op.prob
}

// Apply applies the crossover operation to the selected candidates.
//
// Pairs of candidates are chosen randomly from the selected candidates and
// subjected to crossover to produce a pair of offspring candidates. The
// selected candidates, sel, are the evolved individuals that have survived to
// be eligible to reproduce.
//
// Returns the combined set of evolved offsprings generated by applying
// crossover to the selected candidates.
func (op *Crossover) Apply(sel []framework.Candidate, rng *rand.Rand) []framework.Candidate {
	// Shuffle the collection before applying each operation so that the
	// evolution is not influenced by any ordering artifacts from previous
	// operations.
	selcopy := make([]framework.Candidate, len(sel))
	copy(selcopy, sel)
	framework.ShuffleCandidates(selcopy, rng)

	res := make([]framework.Candidate, 0, len(sel))
	for i := 0; i < len(selcopy); {
		p1 := selcopy[i]
		i++
		if i < len(selcopy) {
			p2 := selcopy[i]
			i++
			// Randomly decide (according to the current crossover probability)
			// whether to perform crossover for these 2 parents.
			var nxpts int64
			if op.prob.NextValue().NextEvent(rng) {
				nxpts = op.npts.NextValue()
			}

			if nxpts > 0 {
				res = append(res, op.Mate(p1, p2, nxpts, rng)...)
			} else {
				// If there is no crossover to perform, just add the parents to the
				// results unaltered.
				res = append(res, p1, p2)
			}
		} else {
			// If we have an odd number of selected candidates, we can't pair up
			// the last one so just leave it unmodified.
			res = append(res, p1)
		}
	}
	return res
}
