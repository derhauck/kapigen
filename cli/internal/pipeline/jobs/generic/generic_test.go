package generic

import (
	"errors"
	"strings"
	"testing"

	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/types"
)

func TestNewGenericJob(t *testing.T) {
	t.Run("valid input parameters", func(t *testing.T) {
		job, err := NewGenericJob("imageName", stages.LINT, []string{})
		if err != nil {
			t.Errorf("Error creating job: %s", err)
		}
		if job == nil {
			t.Error("Job was nil")
		}
	})
	t.Run("invalid input parameters", func(t *testing.T) {
		job, err := NewGenericJob("", stages.LINT, []string{})
		if err == nil {
			t.Error("Error was nil")
		}
		var detailErr *types.DetailedError
		if !errors.As(err, &detailErr) {
			t.Error("Error was not a DetailedError")
		}
		if detailErr.Filename != "generic.go" && !strings.Contains(detailErr.Msg, "imageName") {
			t.Errorf("Error was not the expected one: %s", detailErr.Error())
		}

		if job != nil {
			t.Error("Job was not nil")
		}
	})
}
