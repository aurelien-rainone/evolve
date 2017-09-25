package evolve

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aurelien-rainone/evolve/framework"
)

// Stepper is the interface implemented by objects having a NextEvolutionStep
// method.
type Stepper interface {

	// NextEvolutionStep performs a single step/iteration of the evolutionary process.
	//
	// - evaluatedPopulation is the population at the beginning of the process.
	// - eliteCount is the number of the fittest individuals that must be
	// preserved.
	//
	// Returns the updated population after the evolutionary process has
	// proceeded by one step/iteration.
	NextEvolutionStep(
		evaluatedPopulation framework.EvaluatedPopulation,
		eliteCount int,
		rng *rand.Rand) framework.EvaluatedPopulation
}

// AbstractEvolutionEngine is a base struc for EvolutionEngine implementations.
type AbstractEvolutionEngine struct {
	// A single multi-threaded worker is shared among multiple evolution engine instances.
	concurrentWorker               *fitnessEvaluationPool
	observers                      map[framework.EvolutionObserver]struct{}
	rng                            *rand.Rand
	candidateFactory               framework.CandidateFactory
	fitnessEvaluator               framework.FitnessEvaluator
	singleThreaded                 bool
	satisfiedTerminationConditions []framework.TerminationCondition
	Stepper
}

// NewAbstractEvolutionEngine creates a new evolution engine by specifying the
// various components required by an evolutionary algorithm.
//
// - candidateFactory is the factory used to create the initial population that
// is iteratively evolved.
// - fitnessEvaluator is a function for assigning fitness scores to candidate
// solutions.
// - rng is the source of randomness used by all stochastic processes (including
// evolutionary operators and selection strategies).
func NewAbstractEvolutionEngine(candidateFactory framework.CandidateFactory,
	fitnessEvaluator framework.FitnessEvaluator,
	rng *rand.Rand,
	stepper Stepper) *AbstractEvolutionEngine {

	return &AbstractEvolutionEngine{
		candidateFactory: candidateFactory,
		fitnessEvaluator: fitnessEvaluator,
		rng:              rng,
		observers:        make(map[framework.EvolutionObserver]struct{}),
		Stepper:          stepper,
	}
}

// Evolve executes the evolutionary algorithm until one of the termination
// conditions is met, then return the fittest candidate from the final
// generation.
//
// To return the entire population rather than just the fittest candidate,
// use the EvolvePopulation method instead.
//
// - populationSize is the number of candidate solutions present in the
// population at any point in time.
// - eliteCount is the number of candidates preserved via elitism. In
// elitism, a sub-set of the population with the best fitness scores are
// preserved unchanged in the subsequent generation. Candidate solutions
// that are preserved unchanged through elitism remain eligible for
// selection for breeding the remainder of the next generation. This value
// must be non-negative and less than the population size. A value of zero
// means that no elitism will be applied.
// - conditions is a slice of conditions that may cause the evolution to
// terminate.
//
// Return the fittest solution found by the evolutionary process.
func (e *AbstractEvolutionEngine) Evolve(populationSize, eliteCount int,
	conditions ...framework.TerminationCondition) framework.Candidate {

	return e.EvolveWithSeedCandidates(populationSize,
		eliteCount,
		[]framework.Candidate{},
		conditions...)
}

// EvolveWithSeedCandidates executes the evolutionary algorithm until one of
// the termination conditions is met, then return the fittest candidate from
// the final generation. Provide a set of candidates to seed the starting
// population with.
//
// To return the entire population rather than just the fittest candidate,
// use the EvolvePopulationWithSeedCandidates method instead.
// - populationSize is the number of candidate solutions present in the
// population at any point in time.
// - eliteCount is the number of candidates preserved via elitism. In
// elitism, a sub-set of the population with the best fitness scores are
// preserved unchanged in the subsequent generation. Candidate solutions
// that are preserved unchanged through elitism remain eligible for
// selection for breeding the remainder of the next generation.  This value
// must be non-negative and less than the population size. A value of zero
// means that no elitism will be applied.
// - seedCandidates is a set of candidates to seed the population with. The
// size of this collection must be no greater than the specified population
// size.
// - conditions is a slice of conditions that may cause the evolution to
// terminate.
//
// Returns the fittest solution found by the evolutionary process.
func (e *AbstractEvolutionEngine) EvolveWithSeedCandidates(populationSize, eliteCount int,
	seedCandidates []framework.Candidate,
	conditions ...framework.TerminationCondition) framework.Candidate {

	return e.EvolvePopulationWithSeedCandidates(populationSize,
		eliteCount,
		seedCandidates,
		conditions...)[0].Candidate()
}

// EvolvePopulation executes the evolutionary algorithm until one of the
// termination conditions is met, then return all of the candidates from the
// final generation.
//
// To return just the fittest candidate rather than the entire population,
// use the Evolve method instead.
// - populationSize is the number of candidate solutions present in the
// population at any point in time.
// - eliteCount is the number of candidates preserved via elitism. In
// elitism, a sub-set of the population with the best fitness scores are
// preserved unchanged in the subsequent generation. Candidate solutions
// that are preserved unchanged through elitism remain eligible for
// selection for breeding the remainder of the next generation.  This value
// must be non-negative and less than the population size. A value of zero
// means that no elitism will be applied.
// -  conditions is a slice of conditions that may cause the evolution to
// terminate.
//
// Return the fittest solution found by the evolutionary process.
func (e *AbstractEvolutionEngine) EvolvePopulation(populationSize, eliteCount int,
	conditions ...framework.TerminationCondition) framework.EvaluatedPopulation {

	return e.EvolvePopulationWithSeedCandidates(populationSize,
		eliteCount,
		[]framework.Candidate{},
		conditions...)
}

// EvolvePopulationWithSeedCandidates executes the evolutionary algorithm
// until one of the termination conditions is met, then return all of the
// candidates from the final generation.
//
// To return just the fittest candidate rather than the entire population,
// use the EvolveWithSeedCandidates method instead.
// - populationSize is the number of candidate solutions present in the
// population at any point in time.
// - eliteCount The number of candidates preserved via elitism.  In elitism,
// a sub-set of the population with the best fitness scores are preserved
// unchanged in the subsequent generation.  Candidate solutions that are
// preserved unchanged through elitism remain eligible for selection for
// breeding the remainder of the next generation.  This value must be
// non-negative and less than the population size.  A value of zero means
// that no elitism will be applied.
// - seedCandidates A set of candidates to seed the population with.  The
// size of this collection must be no greater than the specified population
// size.
// - conditions One or more conditions that may cause the evolution to
// terminate.
//
// Return the fittest solution found by the evolutionary process.
func (e *AbstractEvolutionEngine) EvolvePopulationWithSeedCandidates(
	populationSize, eliteCount int,
	seedCandidates []framework.Candidate,
	conditions ...framework.TerminationCondition) framework.EvaluatedPopulation {

	if eliteCount < 0 || eliteCount >= populationSize {
		panic("Elite count must be non-negative and less than population size.")
	}
	if len(conditions) == 0 {
		panic("At least one TerminationCondition must be specified.")
	}

	e.satisfiedTerminationConditions = nil
	var currentGenerationIndex int
	startTime := time.Now()

	population := e.candidateFactory.SeedInitialPopulation(populationSize,
		seedCandidates,
		e.rng)

	// Calculate the fitness scores for each member of the initial population.
	evaluatedPopulation := e.evaluatePopulation(population)

	SortEvaluatedPopulation(evaluatedPopulation, e.fitnessEvaluator.IsNatural())
	data := ComputePopulationData(evaluatedPopulation,
		e.fitnessEvaluator.IsNatural(),
		eliteCount,
		currentGenerationIndex,
		startTime)

	// Notify observers of the state of the population.
	e.notifyPopulationChange(data)

	satisfiedConditions := ShouldContinue(data, conditions...)
	for satisfiedConditions == nil {
		currentGenerationIndex++
		evaluatedPopulation = e.NextEvolutionStep(evaluatedPopulation, eliteCount, e.rng)
		SortEvaluatedPopulation(evaluatedPopulation, e.fitnessEvaluator.IsNatural())
		data = ComputePopulationData(evaluatedPopulation,
			e.fitnessEvaluator.IsNatural(),
			eliteCount,
			currentGenerationIndex,
			startTime)
		// Notify observers of the state of the population.
		e.notifyPopulationChange(data)
		satisfiedConditions = ShouldContinue(data, conditions...)
	}
	e.satisfiedTerminationConditions = satisfiedConditions
	return evaluatedPopulation
}

// Takes a population, assigns a fitness score to each member and returns
// the members with their scores attached, sorted in descending order of
// fitness (descending order of fitness score for natural scores, ascending
// order of scores for non-natural scores).
// - population is the population to evaluate (each candidate is assigned a
// fitness score).
//
// Returns the evaluated population (a list of candidates with attached fitness
// scores).
func (e *AbstractEvolutionEngine) evaluatePopulation(population []framework.Candidate) framework.EvaluatedPopulation {
	// TODO: change comment about thread
	var evaluatedPopulation framework.EvaluatedPopulation
	// Do fitness evaluations on the request thread.
	var err error
	if e.singleThreaded {
		evaluatedPopulation = make(framework.EvaluatedPopulation, len(population))
		for i, candidate := range population {
			evaluatedPopulation[i], err = framework.NewEvaluatedCandidate(candidate, e.fitnessEvaluator.Fitness(candidate, population))
			if err != nil {
				panic(fmt.Sprintf("Can't evaluate candidate %v: %v", candidate, err))
			}
		}
	} else {
		// Divide the required number of fitness evaluations equally among the
		// available processors and coordinate the threads so that we do not
		// proceed until all threads have finished processing.
		unmodifiablePopulation := make([]framework.Candidate, len(population))
		// TODO: is this really necessary?
		copy(unmodifiablePopulation, population)

		// Submit tasks for execution and wait until all threads have finished fitness evaluations.
		evaluatedPopulation = e.pool().submit(
			newFitnessEvaluationTask(
				e.fitnessEvaluator,
				unmodifiablePopulation,
			))

		// TODO: handle goroutine termination
		/*
		   catch (InterruptedException ex)
		   {
		       // Restore the interrupted status, allows methods further up the call-stack
		       // to abort processing if appropriate.
		       Thread.currentThread().interrupt();
		   }
		*/
	}

	return evaluatedPopulation
}

// SatisfiedTerminationConditions returns a slice of all TerminationCondition's
// that are satisfied by the current state of the evolution engine.
//
// Usually this slice will contain only one item, but it is possible that
// mutliple termination conditions will become satisfied at the same time. In
// this case the condition objects in the slice will be in the same order that
// they were specified when passed to the engine.
//
// If the evolution has not yet terminated (either because it is still in
// progress or because it hasn't even been started) then
// framework.ErrIllegalState is returned.
//
// If the evolution terminated because the request thread was interrupted before
// any termination conditions were satisfied then this method will return an
// empty slice.
//
// The slice is guaranteed to be non-null. The slice may be empty because it is
// possible for evolution to terminate without any conditions being matched.
// The only situation in which this occurs is when the request thread is
// interrupted.
func (e *AbstractEvolutionEngine) SatisfiedTerminationConditions() ([]framework.TerminationCondition, error) {
	if e.satisfiedTerminationConditions == nil {
		//throw new IllegalStateException("EvolutionEngine has not terminated.");
		return nil, framework.ErrIllegalState("evolution engine has not terminated")
	}
	satisfiedTerminationConditions := make([]framework.TerminationCondition, len(e.satisfiedTerminationConditions))
	copy(satisfiedTerminationConditions, e.satisfiedTerminationConditions)
	return satisfiedTerminationConditions, nil
}

// AddEvolutionObserver adds a listener to receive status updates on the
// evolution progress.
//
// Updates are dispatched synchronously on the request thread. Observers should
// complete their processing and return in a timely manner to avoid holding up
// the evolution.
func (e *AbstractEvolutionEngine) AddEvolutionObserver(observer framework.EvolutionObserver) {
	e.observers[observer] = struct{}{}
}

// RemoveEvolutionObserver removes an evolution progress listener.
func (e *AbstractEvolutionEngine) RemoveEvolutionObserver(observer framework.EvolutionObserver) {
	delete(e.observers, observer)
}

// notifyPopulationChange sends the population data to all registered observers.
func (e *AbstractEvolutionEngine) notifyPopulationChange(data *framework.PopulationData) {
	for observer := range e.observers {
		observer.PopulationUpdate(data)
	}
}

// SetSingleThreaded forces evaluation to occur synchronously on the request
// goroutine.
//
// By default, fitness evaluations are performed on separate goroutines (as many
// as there are available cores/processors). This is useful in restricted
// environments where programs are not permitted to start or control threads. It
// might also lead to better performance for programs that have extremely
// lightweight/trivial fitness evaluations.
func (e *AbstractEvolutionEngine) SetSingleThreaded(singleThreaded bool) {
	e.singleThreaded = singleThreaded
}

// pool lazily creates the fitness evaluations goroutine pool.
func (e *AbstractEvolutionEngine) pool() *fitnessEvaluationPool {
	if e.concurrentWorker == nil {
		e.concurrentWorker = newFitnessEvaluationPool()
	}
	return e.concurrentWorker
}
