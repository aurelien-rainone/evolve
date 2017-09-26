package framework

import (
	"math/rand"
)

// An EvolutionaryOperator is a function that takes a population of candidates
// as an argument and returns a new population that is the result of applying a
// transformation to the original population.
//
// An implementation of this interface must not modify any of the selected
// candidate objects passed in. Doing so will affect the correct operation of
// the EvolutionEngine. Instead the operator should create and return new
// candidate objects.
//
// However, the operator is not required to create copies of unmodified
// individuals (for efficiency these may be returned directly).
type EvolutionaryOperator interface {

	// Apply applies the operation to each entry in the list of selected
	// candidates.
	//
	// It is important to note that this method operates on the list of
	// candidates returned by the selection strategy and not on the current
	// population. Each entry in the list (not each individual - the list may
	// contain the same individual more than once) must be operated on exactly
	// once.
	//
	// Implementing structs should not assume any particular ordering (or lack
	// of ordering) for the selection. If ordering or shuffling is required, it
	// should be performed by the implementing struct. The implementation should
	// not re-order the list provided but instead should make a copy of the list
	// and re-order that.
	//
	// The ordering of the selection should be totally irrelevant for operators
	// that process each candidate in isolation, such as mutation.  It should
	// only be an issue for operators, such as cross-over, that deal with
	// multiple candidates in a single operation.
	//
	// The operator must not modify any of the candidates passed. Instead it
	// should return a list that contains evolved copies of those candidates
	// (umodified candidates can be included in the results without having to be
	// copied).
	Apply(selectedCandidates []Candidate, rng *rand.Rand) []Candidate
}