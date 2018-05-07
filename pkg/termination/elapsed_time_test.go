package termination

import (
	"testing"
	"time"

	"github.com/aurelien-rainone/evolve/pkg/api"
)

func TestElapsedTime(t *testing.T) {
	cond := 1 * ElapsedTime(time.Second)
	popdata := &api.PopulationData{}

	popdata.Elapsed = 100 * time.Millisecond
	if cond.ShouldTerminate(popdata) {
		t.Errorf("should not terminate before elapsed time")
	}

	popdata.Elapsed = time.Second
	if !cond.ShouldTerminate(popdata) {
		t.Errorf("should terminate after elapsed time")
	}
}