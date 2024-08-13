package job

import (
	"reflect"
	"testing"
)

func TestEnvironmentActionEnum(t *testing.T) {
	t.Run("can run enum", func(t *testing.T) {
		values := EnvironmentActionEnum().GetValues()
		if len(values) == 0 {
			t.Error("should not be empty")
			t.FailNow()
		}
	})
	t.Run("can not get enum", func(t *testing.T) {
		_, err := EnvironmentActionEnum().Value(99)
		if err == nil {
			t.Error("should have error")
			t.FailNow()
		}
	})
	t.Run("can get string", func(t *testing.T) {
		if EnvironmentActionStart.String() != "start" {
			t.Error("should be start, received: ", EnvironmentActionStart.String())
			t.FailNow()
		}
	})
}
func TestEnvironmentTierEnum(t *testing.T) {
	t.Run("can run enum", func(t *testing.T) {
		values := EnvironmentTierEnum().GetValues()
		if len(values) == 0 {
			t.Error("should have error")
			t.FailNow()
		}
	})

	t.Run("can not get enum", func(t *testing.T) {
		_, err := EnvironmentTierEnum().Value(99)
		if err == nil {
			t.Error("should not be empty")
			t.FailNow()
		}
	})

	t.Run("can get string", func(t *testing.T) {
		if EnvironmentTierDeployment.String() != "development" {
			t.Error("should be development, received: ", EnvironmentTierDeployment.String())
			t.FailNow()
		}
	})
}

func TestEnvironment_Render(t *testing.T) {

	expected := &EnvironmentYaml{
		DeploymentTier: EnvironmentTierEnum().ValueSafe(EnvironmentTierProduction),
		Action:         EnvironmentActionEnum().ValueSafe(EnvironmentActionStart),
		Name:           "test",
		Url:            "https://test.url",
		AutoStopIn:     "2d",
	}

	actual := &Environment{
		Name:           "test",
		Url:            "https://test.url",
		Action:         EnvironmentActionStart,
		AutoStopIn:     "2d",
		DeploymentTier: EnvironmentTierProduction,
	}
	if !reflect.DeepEqual(actual.Render(), expected) {
		t.Errorf("expected %v, received %v", expected, actual.Render())
	}
}
