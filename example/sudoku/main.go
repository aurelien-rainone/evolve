package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/arl/evolve/pkg/api"
	"github.com/arl/evolve/pkg/engine"
	"github.com/arl/evolve/pkg/operator"
	"github.com/arl/evolve/pkg/operator/xover"
	"github.com/arl/evolve/pkg/selection"
	"github.com/arl/evolve/pkg/termination"
	"github.com/arl/evolve/random"
)

func check(err error, v ...interface{}) {
	if err != nil {
		if len(v) == 0 {
			log.Fatal(v, err)
		}
	}
}

func readSudokus(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close() // nolint: errcheck
	}()

	names, err := f.Readdirnames(0)
	switch {
	case err != nil:
		return nil, err
	case len(names) == 0:
		return nil, errors.New("empty directory")
	}
	return names, err
}

func readPattern(fn string) ([]string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close() // nolint: errcheck
	}()

	puzzle := []string{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		puzzle = append(puzzle, s.Text())
	}
	return puzzle, s.Err()
}

func solveSudoku(pattern []string) error {
	rng := rand.New(random.NewMT19937(time.Now().UnixNano()))

	// Crossover rows between parents (so offspring is x rows from parent1 and y
	// rows from parent2).
	xover := xover.New(mater{})
	check(xover.SetPoints(1))

	mutation := newRowMutation()
	// TODO: use a PoissonGenerator for mutation count and a
	// DiscreteUniformGenerator for mutation amount
	check(mutation.SetMutationsRange(1, 2))
	check(mutation.SetAmountRange(1, 8))

	pipeline := operator.Pipeline{xover, mutation}

	selector := selection.NewTournament()
	check(selector.SetProb(0.85))

	gen, err := newGenerator(pattern)
	check(err)

	eng := engine.NewGenerational(gen, pipeline, evaluator{}, selector, rng)

	eng.AddObserver(api.ObserverFunc(func(data *api.PopulationData) {
		if data.GenNumber%100 == 0 {
			return
		}
		// only shows multiple of 100 generations
		fmt.Printf("At generation %d the best solution has a fitness of %v\n%v\n", data.GenNumber, data.BestFitness, data.BestCand.(*sudoku))
	}))

	const (
		popsize = 500
		nelites = 500 * 0.05
	)
	solution := eng.Evolve(
		popsize, nelites,
		termination.TargetFitness{Fitness: 0, Natural: false},
		termination.NewUserAbort(),
	)

	fmt.Printf("solution:\n%v\n", solution.(*sudoku))
	return nil
}

func main() {
	puzdir := flag.String("puzzles", "./puzzles", "directory containing the puzzles")
	flag.Parse()

	puzzles, err := readSudokus(*puzdir)
	if err != nil {
		log.Fatalf("can't read puzzles directory: %v", err)
	}

	for i, p := range puzzles {
		fmt.Printf("\t[%d] %s\n", i, p)
	}

	fmt.Print("Choose the sudoku puzzle you want to solve? ")
	var i int
	if _, err = fmt.Scanf("%d", &i); err != nil {
		log.Fatalf("can't read your choice: %v", err)
		return
	}
	if i < 0 || i >= len(puzzles) {
		log.Fatal("invalid entry")
	}

	pattern, err := readPattern(path.Join(*puzdir, puzzles[i]))
	if err != nil {
		log.Fatalf("can't read sudo pattern: %v", err)
	}

	err = solveSudoku(pattern)
	if err != nil {
		log.Fatalf("couldn't solve sudoku pattern: %v\n", err)
	}
}
