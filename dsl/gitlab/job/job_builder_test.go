package job

import (
	"reflect"
	"testing"
)

func TestCiJob_AddSecret(t *testing.T) {
	t.Run("can add from NewVaultSecret", func(t *testing.T) {
		expected := &CiJob{
			Secrets: Secrets{
				"TEST": &VaultSecret{
					Vault: VaultSecretConfig{
						Engine: VaultSecretEngine{
							Name: EnumVaultSecretEngineKv2,
							Path: "mount",
						},
						Path:  "path",
						Field: "field",
					},
					Token: "token",
				},
			},
		}
		job := &CiJob{}
		job.AddSecret("TEST", NewVaultSecret(EnumVaultSecretEngineKv2, "mount", "path", "field", "token"))
		actual := job
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})
	t.Run("can add from NewVaultKv2Secret", func(t *testing.T) {
		expected := &CiJob{
			Secrets: Secrets{
				"TEST": &VaultSecret{
					Vault: VaultSecretConfig{
						Engine: VaultSecretEngine{
							Name: EnumVaultSecretEngineKv2,
							Path: "mount",
						},
						Path:  "path",
						Field: "field",
					},
					Token: "token",
				},
			},
		}
		job := &CiJob{}
		job.AddSecret("TEST", NewVaultKv2Secret("mount", "path", "field", "token"))
		actual := job
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})
	t.Run("can overwrite already present secret", func(t *testing.T) {
		job := &CiJob{
			Secrets: Secrets{
				"TEST": &VaultSecret{
					Vault: VaultSecretConfig{
						Engine: VaultSecretEngine{
							Name: EnumVaultSecretEngineKv2,
							Path: "mount",
						},
						Path:  "path",
						Field: "field",
					},
					Token: "token",
				},
			},
		}
		expected := NewVaultSecret(EnumVaultSecretEngineKv1, "mount_fail", "path_fail", "field_fail", "token_fail")
		job.AddSecret("TEST", expected)
		actual := job.Secrets["TEST"]
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", job, actual)
		}
	})
}

func TestCiJob_AddIdToken(t *testing.T) {
	t.Run("can add new IdToken", func(t *testing.T) {
		job := &CiJob{}
		job.AddIdToken("TEST", "https://test.aud")

		if job.IdTokens["TEST"].AUD != "https://test.aud" {
			t.Errorf("expected %v, got %v", "https://test.aud", job.IdTokens["TEST"])
		}
	})

	t.Run("can overwrite IdToken", func(t *testing.T) {

		job := &CiJob{
			IdTokens: IdTokens{
				"TEST": &IdToken{
					AUD: "https://unexpected.aud",
				},
			},
		}
		job.AddIdToken("TEST", "https://test.aud")

		if job.IdTokens["TEST"].AUD != "https://test.aud" {
			t.Errorf("expected %v, got %v", "https://test.aud", job.IdTokens["TEST"])
		}
	})
}
