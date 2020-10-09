package main

import (
	"github.com/arl/evolve"
	"github.com/arl/evolve/engine"
	"github.com/arl/evolve/generator"
	"github.com/arl/evolve/operator"
	"github.com/arl/evolve/operator/xover"
	"github.com/arl/evolve/pkg/bitstring"
	"github.com/arl/evolve/selection"
)

func runTSP() error {
	// Define the crossover operator.
	xover := xover.New(xover.PMX)
	xover.Points = generator.ConstInt(1)
	xover.Probability = generator.ConstFloat64(0.7)

	// Define the mutation operator.
	mut := mutation.NewBitstring()
	op := ListOrder{
		Count:  generator.NewPoisson(generator.ConstFloat64(1.5)),
		Amount: generator.NewPoisson(generator.ConstFloat64(1.5)),
	}

	eval := evolve.EvaluatorFunc(
		true, // natural fitness (higher is better)
		func(cand interface{}, pop []interface{}) float64 {
			// our evaluator counts the ones in the bitstring
			return float64(cand.(*bitstring.Bitstring).OnesCount())
		})

	epocher := engine.Generational{
		Op:   operator.Pipeline{xover, mut},
		Eval: eval,
		Sel:  selection.RouletteWheel,
	}

	eng, err := engine.New(generator.Bitstring(nbits), eval, &epocher)
	check(err)
}

/*
{
        Random rng = new MersenneTwisterRNG();

        // Set-up evolution pipeline (cross-over followed by mutation).
        List<EvolutionaryOperator<List<String>>> operators = new ArrayList<EvolutionaryOperator<List<String>>>(2);
        if (crossover)
        {
            operators.add(new ListOrderCrossover<String>());
        }
        if (mutation)
        {
            operators.add(new ListOrderMutation<String>(new PoissonGenerator(1.5, rng),
                                                        new PoissonGenerator(1.5, rng)));
        }

        EvolutionaryOperator<List<String>> pipeline = new EvolutionPipeline<List<String>>(operators);

        CandidateFactory<List<String>> candidateFactory
            = new ListPermutationFactory<String>(new LinkedList<String>(cities));
        EvolutionEngine<List<String>> engine
            = new GenerationalEvolutionEngine<List<String>>(candidateFactory,
                                                            pipeline,
                                                            new RouteEvaluator(distances),
                                                            selectionStrategy,
                                                            rng);
        if (progressListener != null)
        {
            engine.addEvolutionObserver(new EvolutionObserver<List<String>>()
            {
                public void populationUpdate(PopulationData<? extends List<String>> data)
                {
                    progressListener.updateProgress(((double) data.getGenerationNumber() + 1) / generationCount * 100);
                }
            });
        }
        return engine.evolve(populationSize,
                             eliteCount,
                             new GenerationCount(generationCount));
    }*/
