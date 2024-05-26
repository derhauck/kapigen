package php

import "testing"

func TestNewPhpUnit(t *testing.T) {
	t.Run("Can create phpunit job", func(t *testing.T) {
		t.Parallel()
		t.Run("Correct parameters", func(t *testing.T) {
			job, err := NewPhpUnit("testimage", "testpath", "test", ".", "", map[string]int32{})
			if err != nil {
				t.Error(err)
			}
			err = job.Render()
			if err != nil {
				t.Error(err)
			}
			if job.CiJobYaml.String() == "" {
				t.Error("Should not be empty")
			}
		})
		t.Run("Incorrect parameters", func(t *testing.T) {
			job, err := NewPhpUnit("", "", "", "", "", map[string]int32{})
			if err != nil {
				t.Error(err)
			}
			err = job.Render()
			if err == nil {
				t.Errorf("created job succeeded without image and path but should not: %s", job.CiJobYaml.String())
			}

		})
	})
}
