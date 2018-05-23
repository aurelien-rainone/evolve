package generator

import (
	"bytes"
	"errors"
	"math/rand"
	"unicode/utf8"
)

var (
	// ErrEmptyAlphabet is the error returned by NewString when providing an
	// empty alphabet string.
	ErrEmptyAlphabet = errors.New("alphabet should not be empty")

	// ErrNotASCIIAlphabet is the error returned by NewString when the alphabet
	// contains some non-ASCII runes.
	ErrNotASCIIAlphabet = errors.New("alphabet should only contain ASCII runes")
)

// String is a generator of ASCII string candidates of the specified length and
// in which runes are randomly chosen from an alphabet
type String struct {
	alphabet string
	length   int
}

// NewString returns a String that generates strings of the specified length
// from the provided alphabet.
//
// NewString will return ErrEmptyAlphabet if the alphabet is empty and
// ErrNotASCIIAlphabet if the alphabet contains some non-ASCII runes.
func NewString(alphabet string, length int) (*String, error) {
	if alphabet == "" {
		return nil, ErrEmptyAlphabet
	}
	if utf8.RuneCountInString(alphabet) != len(alphabet) {
		return nil, ErrNotASCIIAlphabet
	}

	return &String{
		alphabet: alphabet,
		length:   length,
	}, nil
}

// GenerateCandidate generates a random string.
func (gen *String) GenerateCandidate(rng *rand.Rand) interface{} {
	var buffer bytes.Buffer
	for i := 0; i < gen.length; i++ {
		idx := rand.Int31n(int32(len(gen.alphabet)))
		buffer.WriteByte(gen.alphabet[idx])
	}
	return buffer.String()
}
