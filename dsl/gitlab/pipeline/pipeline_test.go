package pipeline

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestNewDefaultCiPipeline(t *testing.T) {
	t.Run("create default pipeline", func(t *testing.T) {
		ciPipeline := (&CiPipeline{}).DefaultCiPipeline()
		stages := ciPipeline.Stages.Get()
		snaps.MatchSnapshot(t, ciPipeline.Default, ciPipeline.Workflow, ciPipeline.Variables, len(stages))
	})

}
