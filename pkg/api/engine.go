package api

// Engine is the interface implemented by objects that provide evolution
// operations.
// TODO: does the Engine interface really needs all this methods? wouldn't one
// suffice and the others be derived from it (in engine.Base)?
// TODO: Could AddObserver/RemoveObserver be made an external interace, included
// in Engine? How would that go with future island observers?
type Engine interface {

	// Evolve executes the evolutionary algorithm until one of the termination
	// conditions is met, then return the fittest candidate from the final
	// generation.
	//
	// To return the entire population rather than just the fittest candidate,
	// use the EvolvePopulation method instead.
	//
	// size is the number of candidate solutions present in the population at
	// any point in time.
	// nelites is the number of candidates preserved via elitism. In elitism, a
	// sub-set of the population with the best fitness scores are preserved
	// unchanged in the subsequent generation. Candidate solutions that are
	// preserved unchanged through elitism remain eligible for selection for
	// breeding the remainder of the next generation. This value must be
	// non-negative and less than the population size. A value of zero means
	// that no elitism will be applied.
	// conds is a slice of conditions that may cause the evolution to terminate.
	//
	// Returns the fittest solution found by the evolutionary process.
	Evolve(size, nelites int, conds ...TerminationCondition) interface{}

	// EvolveWithSeedCandidates executes the evolutionary algorithm until one of
	// the termination conditions is met, then return the fittest candidate from
	// the final generation. Provide a set of candidates to seed the starting
	// population with.
	//
	// To return the entire population rather than just the fittest candidate,
	// use the EvolvePopulationWithSeedCandidates method instead.
	//
	// size is the number of candidate solutions present in the population at
	// any point in time.
	// nelites is the number of candidates preserved via elitism. In elitism, a
	// sub-set of the population with the best fitness scores are preserved
	// unchanged in the subsequent generation. Candidate solutions that are
	// preserved unchanged through elitism remain eligible for selection for
	// breeding the remainder of the next generation.  This value must be
	// non-negative and less than the population size. A value of zero means
	// that no elitism will be applied.
	// seedcands is a set of candidates to seed the population with. The size of
	// this collection must be no greater than the specified population size.
	// conds is a slice of conditions that may cause the evolution to terminate.
	//
	// Returns the fittest solution found by the evolutionary process.
	EvolveWithSeedCandidates(size, nelites int, seedcands []interface{},
		conds ...TerminationCondition) interface{}

	// EvolvePopulation executes the evolutionary algorithm until one of the
	// termination conditions is met, then return all of the candidates from the
	// final generation.
	//
	// To return just the fittest candidate rather than the entire population,
	// use the Evolve method instead.
	// size is the number of candidate solutions present in the population at
	// any point in time.
	// nelites is the number of candidates preserved via elitism. In elitism, a
	// sub-set of the population with the best fitness scores are preserved
	// unchanged in the subsequent generation. Candidate solutions that are
	// preserved unchanged through elitism remain eligible for selection for
	// breeding the remainder of the next generation.  This value must be
	// non-negative and less than the population size. A value of zero means
	// that no elitism will be applied.
	// conds is a slice of conditions that may cause the evolution to terminate.
	//
	// Returns the fittest solution found by the evolutionary process.
	EvolvePopulation(size, nelites int, conds ...TerminationCondition) Population

	// EvolvePopulationWithSeedCandidates executes the evolutionary algorithm
	// until one of the termination conditions is met, then return all of the
	// candidates from the final generation.
	//
	// To return just the fittest candidate rather than the entire population,
	// use the EvolveWithSeedCandidates method instead.
	// size is the number of candidate solutions present in the population at
	// any point in time.
	// nelites The number of candidates preserved via elitism. In elitism, a
	// sub-set of the population with the best fitness scores are preserved
	// unchanged in the subsequent generation. Candidate solutions that are
	// preserved unchanged through elitism remain eligible for selection for
	// breeding the remainder of the next generation. This value must be
	// non-negative and less than the population size. A value of zero means
	// that no elitism will be applied.
	// seedcands is a set of candidates to seed the population with.The size of
	// this collection must be no greater than the specified population size.
	// conditions One or more conditions that may cause the evolution to
	// terminate.
	//
	// Returns the fittest solution found by the evolutionary process.
	EvolvePopulationWithSeedCandidates(size, nelites int, seedcands []interface{},
		conds ...TerminationCondition) Population

	// AddObserver registers an observer to receive status updates on the
	// evolution progress.
	AddObserver(o Observer)

	// RemoveObserver removes an evolution observer.
	RemoveObserver(o Observer)

	// SatisfiedTerminationConditions returns a slice of all
	// TerminationCondition's that are satisfied by the current state of the
	// evolution engine.
	//
	// Usually this list will contain only one item, but it is possible that
	// multiple termination conditions will become satisfied at the same time.
	// In this case the condition objects in the list will be in the same order
	// that they were specified when passed to the engine.
	//
	// If the evolution has not yet terminated (either because it is still in
	// progress or because it hasn't even been started) then an
	// IllegalStateException will be thrown.
	//
	// If the evolution terminated because the request thread was interrupted
	// before any termination conditions were satisfied then this method will
	// return an empty list.
	//
	// Returns a list of statisfied conditions. The list is guaranteed to be
	// non-null. The list may be empty because it is possible for evolution to
	// terminate without any conditions being matched. The only situation in
	// which this occurs is when the request goroutine is interrupted.
	//
	// May return ErrIllegalState if this method is invoked on an evolution
	// engine before evolution is started or while it is still in progress.
	// TODO: find shorter name 'SatisfiedConditions' ?
	SatisfiedTerminationConditions() ([]TerminationCondition, error)
}