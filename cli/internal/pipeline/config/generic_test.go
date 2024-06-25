package config

import (
	"reflect"
	"testing"

	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/cli"
	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/version"
)

func TestGeneric_Build(t *testing.T) {
	t.Run("can build a generic pipeline", func(t *testing.T) {
		expectedScripts := []string{"ls", "pwd"}
		generic := &Generic{
			Stage:   stages.TEST.String(),
			Scripts: expectedScripts,
		}
		main := factory.New(&cli.Settings{
			Mode:         version.Gitlab,
			PrivateToken: "",
		})
		err := generic.Validate()
		jobs, err := generic.Build(main, "generic", "generic")
		if err != nil {
			t.Error(err)
		}
		if jobs == nil {
			t.Error("jobs is nil")
			t.Fail()
		}
		genericJob := jobs.GetJobs()[0]
		if genericJob == nil {
			t.Error("genericJob is nil")
			t.Fail()
		}
		err = genericJob.Render()
		if err != nil {
			t.Error(err)
		}
		if genericJob.CiJob.Image.Name != docker.Alpine_3_18.String() {
			t.Errorf("expected image name to be %s, received %s", docker.Alpine_3_18.String(), genericJob.CiJob.Image.Name)
		}
		if !reflect.DeepEqual(genericJob.CiJobYaml.Script, expectedScripts) {
			t.Errorf("expected scripts to be %s, received %s", expectedScripts, genericJob.CiJobYaml.Script)
		}
	})

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
