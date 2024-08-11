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

func TestVaultSecretYaml_SecretYaml(t *testing.T) {
	t.Run("can get VaultSecretYaml from SecretYaml function", func(t *testing.T) {
		yaml := VaultSecretYaml{
			Vault: VaultSecretConfigYaml{
				Engine: VaultSecretEngineYaml{
					Path: "mount",
					Name: "kv-v2",
				},
				Path:  "path",
				Field: "field",
			},
		}

		result := yaml.SecretYaml()
		switch secret := result.(type) {
		case *VaultSecretYaml:
			t.Log("success")
			if secret.Vault.Engine.Path != yaml.Vault.Engine.Path {
				t.Errorf("expected %v, received %v", yaml.Vault.Engine.Path, secret.Vault.Engine.Path)
			}

			if secret.Vault.Engine.Name != yaml.Vault.Engine.Name {
				t.Errorf("expected %v, received %v", yaml.Vault.Engine.Name, secret.Vault.Engine.Name)
			}

			if secret.Vault.Path != yaml.Vault.Path {
				t.Errorf("expected %v, received %v", yaml.Vault.Path, secret.Vault.Path)
			}

			if secret.Vault.Field != yaml.Vault.Field {
				t.Errorf("expected %v, received %v", yaml.Vault.Field, secret.Vault.Field)
			}

		default:
			t.Errorf("expected *VaultSecretYaml, received %v", result)
		}

	})
}
