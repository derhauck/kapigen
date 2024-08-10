package job

import (
	"reflect"
	"testing"
)

func TestEnumVaultEngineName(t *testing.T) {
	t.Run("can use enum", func(t *testing.T) {
		value, err := EnumVaultSecretEngineName().Value(EnumVaultSecretEngineKv1)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		if value != "kv-v1" {
			t.Errorf("error")
		}
	})

	t.Run("can not use own value", func(t *testing.T) {
		const ENUMTEST VaultSecretEngineName = 75
		value, err := EnumVaultSecretEngineName().Value(ENUMTEST)
		if err == nil {
			t.Error("Expecting Error as value is not present in enum")
			t.FailNow()
		}
		if value != "" {
			t.Errorf("expected value to be empty, received: %s", value)
		}

	})
}

func TestNewVaultSecretYaml(t *testing.T) {
	t.Run("can create valid new secret yaml", func(t *testing.T) {
		expectation := &VaultSecretYaml{
			Vault: VaultSecretConfigYaml{
				Engine: VaultSecretEngineYaml{
					Path: "mount",
					Name: "kv-v2",
				},
				Path:  "path",
				Field: "field",
			},

			Token: "token",
		}
		yaml := NewVaultSecretYaml(&VaultSecret{
			Vault: VaultSecretConfig{
				Engine: VaultSecretEngine{
					Path: "mount",
					Name: EnumVaultSecretEngineKv2,
				},
				Path:  "path",
				Field: "field",
			},

			Token: "token",
		})

		if !reflect.DeepEqual(yaml, expectation) {
			t.Errorf("expected %v, received %v", expectation, yaml)
		}
	})
}
