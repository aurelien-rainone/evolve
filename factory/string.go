package factory

import (
	"errors"
	"math/rand"
	"unicode"
)

var (
	errEmptyAlphabet    = errors.New("alphabet can't be empty")
	errNotASCIIAlphabet = errors.New("alphabet can only contain ASCII runes")
)

// String is a factory for creating random strings with chars taken from an
// alphabet of ASCII chars.
type String struct {
	alphabet []byte
	length   int
}

// NewString returns a factory creating random strings of the specified length
// where chars are taken from the given alphabet.
//
// NewString returns an error if alphabet is empty or contains non-ASCII
// characters.
func NewString(alphabet string, length int) (*String, error) {
	if alphabet == "" {
		return nil, errEmptyAlphabet
	}

	if !isASCII(alphabet) {
		return nil, errNotASCIIAlphabet
	}

	return &String{
		alphabet: []byte(alphabet),
		length:   length,
	}, nil
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// New creates a random string.
func (gen *String) New(rng *rand.Rand) interface{} {
	b := make([]byte, gen.length)
	for i := 0; i < gen.length; i++ {
		b[i] = gen.alphabet[rand.Int31n(int32(len(gen.alphabet)))]
	}
	return string(b)
}
