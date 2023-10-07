package config

import "testing"

func TestInfrastructure_New(t *testing.T) {
	t.Run("Can create new Infrastructure", func(t *testing.T) {
		t.Parallel()
		oldInfra := Infrastructure{
			State: "not-default",
		}
		newInfra := oldInfra.New()
		if result, ok := newInfra.(*Infrastructure); ok {
			if oldInfra.State == result.State {
				t.Error("Could not create NEW struct from type Docker")
			}
		} else {
			t.Error("Not type of Infrastructure")
		}
	})
}
