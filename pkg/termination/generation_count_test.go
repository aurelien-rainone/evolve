package termination

import (
	"testing"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

func TestGenerationCount(t *testing.T) {
	cond := GenerationCount(5)
	popdata := &api.PopulationData{}

	popdata.GenNumber = 3
	if cond.ShouldTerminate(popdata) {
		t.Errorf("should not terminate after 4th generation")
	}

	popdata.GenNumber = 4
	if !cond.ShouldTerminate(popdata) {
		t.Errorf("should terminate after 5th generation")
	}
}
