package php

import "testing"

func TestNewPhpUnit(t *testing.T) {
	t.Run("Can create phpunit job", func(t *testing.T) {
		t.Parallel()
		t.Run("Correct parameters", func(t *testing.T) {
			job, err := NewPhpUnit("testimage", "testpath", "test", ".", "", "testpath", map[string]int32{})
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
			_, err := NewPhpUnit("", "", "", "", "", "testpath", map[string]int32{})
			if err == nil {
				t.Error("Should be missing args error for 'imageName'")
			}
			_, err = NewPhpUnit("testimage", "", "", "", "", "testpath", map[string]int32{})
			if err == nil {
				t.Error("Should be missing args error for 'composerPath'")
			}

		})
	})
}
