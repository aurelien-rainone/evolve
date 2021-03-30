package main

import (
	"fmt"
	"math/rand"
)

var (
	errWrongNumberOfRows     = fmt.Errorf("sudoku layout must have %v rows", size)
	errWrongNumberOfCols     = fmt.Errorf("sudoku layout must have %v cells in each row", size)
	errPatternUnexpectedChar = fmt.Errorf("unexpected char in pattern")
	values                   = [size]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
)

// SudokuGenerator that generates potential Sudoku solutions from a list of "givens".
// The rows of the generated solutions will all be valid (i.e. no duplicate
// values) but there are no constraints on the columns or sub-grids (these will
// be refined by the evolutionary algorithm).
type SudokuGenerator struct {
	templ    sudoku
	nonfixed [size][]int
}

// Creates a generator of random candidate solutions for a specified Sudoku
// puzzle. pattern is a slice of strings, each representing one row of sudoku.
// Each character represents a single cell. Permitted characters are the digits
// '1' to '9' (each of which represents a fixed cell in the pattern) or the '.'
// character, which represents an empty cell. Returns an error if the pattern is
// not made of 9 strings containing 1 to 9, or '.'
func newGenerator(pattern []string) (*SudokuGenerator, error) { // nolint: gocyclo
	if len(pattern) != size {
		return nil, errWrongNumberOfRows
	}

	gen := &SudokuGenerator{}

	// Keep track of which values in each row are not 'givens'.
	for i := 0; i < size; i++ {
		gen.nonfixed[i] = make([]int, size)
		copy(gen.nonfixed[i], values[:])
	}

	for i := 0; i < len(pattern); i++ {
		prow := []byte(pattern[i])

		if len(prow) != size {
			return nil, errWrongNumberOfCols
		}
		for j := 0; j < len(prow); j++ {
			c := prow[j]
			switch {
			case c >= '1' && c <= '9':
				// cell is a 'given'.
				val := int(c - '0')
				gen.templ[i][j].val = val
				gen.templ[i][j].fixed = true

				for idx, cval := range gen.nonfixed[i] {
					if val == cval {
						gen.nonfixed[i] = append(gen.nonfixed[i][:idx], gen.nonfixed[i][idx+1:]...)
						break
					}
				}
			case c == '.':
			default:
				return nil, errPatternUnexpectedChar
			}
		}
	}
	return gen, nil
}

// Potential solutions created by New are guaranteed to have no duplicates in
// any row but could have duplicates in a column or sub-grid.
func (gen *SudokuGenerator) New(rng *rand.Rand) interface{} {
	// Clone the template as the basis for this grid.
	var rows sudoku
	copy(rows[:], gen.templ[:])

	// Fill-in the non-fixed cells.
	for i := 0; i < len(rows); i++ {
		rowvals := gen.nonfixed[i]

		rng.Shuffle(len(rowvals), func(i, j int) {
			rowvals[i], rowvals[j] = rowvals[j], rowvals[i]
		})

		var idx int
		for j := 0; j < len(rows[i]); j++ {
			// because the zero value is invalid, valid values are in [1,9]
			if rows[i][j].val == 0 {
				rows[i][j].val = rowvals[idx]
				rows[i][j].fixed = false
				idx++
			}
		}
	}
	return &rows
}
