package config

import (
	"testing"
)

func TestGeneric_Build(t *testing.T) {

}

func TestGeneric_New(t *testing.T) {
	t.Run("can create a new Generic", func(t *testing.T) {
		old := &Generic{}
		generic := old.New()
		if generic == nil {
			t.Error("Generic was nil")
		}
		if generic == old {
			t.Error("Generic was not a new instance")
		}

		if _, ok := generic.(*Generic); !ok {
			t.Error("returned instance was not of type Generic")
		}
	})
}

func TestGeneric_Rules(t *testing.T) {

}

func TestGeneric_Validate(t *testing.T) {
	t.Run("valid input parameters", func(t *testing.T) {
		generic := &Generic{}
		err := generic.Validate()
		if err != nil {
			t.Errorf("Error validating Generic: %s", err)
		}
	})
}
