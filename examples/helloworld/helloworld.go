package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/aurelien-rainone/evolve"
	"github.com/aurelien-rainone/evolve/factory"
	"github.com/aurelien-rainone/evolve/framework"
	"github.com/aurelien-rainone/evolve/number"
	"github.com/aurelien-rainone/evolve/operators"
	"github.com/aurelien-rainone/evolve/selection"
	"github.com/aurelien-rainone/evolve/termination"
)

func check(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

type observer struct{}

func (o observer) PopulationUpdate(data *framework.PopulationData) {
	fmt.Printf("Generation %d: %s (%v)\n", data.GenerationNumber(), data.BestCandidate(),
		data.BestCandidateFitness())
}

func main() {
	var targetString = "HELLO WORLD"
	if len(os.Args) == 2 {
		targetString = strings.ToUpper(os.Args[1])
	}

	// Create a factory to generate random 11-character Strings.
	alphabet := make([]byte, 27)
	for c := byte('A'); c <= 'Z'; c++ {
		alphabet[c-'A'] = c
	}
	alphabet[26] = ' '

	for _, c := range targetString {
		if !strings.ContainsRune(string(alphabet), c) {
			fmt.Printf("Rune %c is not contained in the alphabet\n", c)
			os.Exit(1)
		}
	}

	var (
		stringFactory *factory.StringFactory
		err           error
	)
	stringFactory, err = factory.NewStringFactory(string(alphabet), len(targetString))
	check(err)

	var (
		mutationProb number.Probability
		mutation     framework.EvolutionaryOperator
		crossover    framework.EvolutionaryOperator
		pipeline     *operators.EvolutionPipeline
	)

	// 1st operator: string mutation
	mutationProb, err = number.NewProbability(0.02)
	check(err)
	mutation, err = operators.NewStringMutation(
		string(alphabet),
		operators.WithConstantStringMutationProbability(mutationProb),
	)
	check(err)

	// 2nd operator: string crossover
	crossover, err = operators.NewStringCrossover()
	check(err)

	// Create a pipeline that applies mutation then crossover
	pipeline, err = operators.NewEvolutionPipeline(mutation, crossover)
	check(err)

	fitnessEvaluator := newStringEvaluator(targetString)

	var selectionStrategy = &selection.RouletteWheelSelection{}
	rng := rand.New(rand.NewSource(99))

	var engine *evolve.AbstractEvolutionEngine
	engine = evolve.NewGenerationalEvolutionEngine(stringFactory,
		pipeline,
		fitnessEvaluator,
		selectionStrategy,
		rng)

	//engine.SetSingleThreaded(true)
	engine.AddEvolutionObserver(observer{})
	result := engine.Evolve(100, 5, termination.NewTargetFitness(0, false))
	fmt.Println(result)

	var conditions []framework.TerminationCondition
	conditions, err = engine.SatisfiedTerminationConditions()
	check(err)
	for i, condition := range conditions {
		fmt.Printf("satified termination condition %v %T: %v\n",
			i, condition, condition)
	}
}