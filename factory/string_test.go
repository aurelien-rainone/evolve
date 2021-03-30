package factory

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewString(t *testing.T) {
	tests := []struct {
		name     string
		alphabet string
		length   int
		wantErr  error
	}{
		{
			name:     "valid string generator",
			alphabet: "abcdefgh12324;?:",
			length:   9,
			wantErr:  nil,
		},
		{
			name:     "empty string generator",
			alphabet: "abcdefgh12324;?:",
			length:   0,
			wantErr:  nil,
		},
		{
			name:     "not ASCII-only alphabet",
			alphabet: "abcdefgh12324本;?:",
			length:   0,
			wantErr:  errNotASCIIAlphabet,
		},
		{
			name:     "empty alphabet",
			alphabet: "",
			length:   10,
			wantErr:  errEmptyAlphabet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewString(tt.alphabet, tt.length)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestStringFactory(t *testing.T) {
	factory, err := NewString("ABCdefg", 2)
	require.NoError(t, err)

	s := factory.New(rand.New(rand.NewSource(99)))
	if s, ok := s.(string); !ok {
		t.Errorf("New should generate string candidates, got %T", s)
	}
}

var sink interface{}

func BenchmarkNewString(b *testing.B) {
	rng := rand.New(rand.NewSource(99))

	runs := []int{10, 100, 1000}
	for _, slen := range runs {
		b.Run(fmt.Sprintf("%d", slen), func(b *testing.B) {
			b.ReportAllocs()
			factory, _ := NewString("A", slen)
			for i := 0; i < b.N; i++ {
				sink = factory.New(rng)
			}
		})
	}
}
