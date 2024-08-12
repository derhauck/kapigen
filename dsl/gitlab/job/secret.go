package job

import (
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

type IdToken struct {
	AUD string
}

type IdTokens map[string]*IdToken

func (i *IdTokens) Render() *IdTokensYaml {
	idTokensYaml := IdTokensYaml{}
	for k, v := range *i {
		idTokensYaml[k] = v.Render()
	}
	return &idTokensYaml
}

type IdTokenYaml struct {
	AUD string `yaml:"aud"`
}

func (i *IdToken) Render() IdTokenYaml {
	return IdTokenYaml{
		AUD: i.AUD,
	}
}

type IdTokensYaml map[string]IdTokenYaml

type VaultSecretEngineName byte

const (
	EnumVaultSecretEngineKv2 VaultSecretEngineName = iota
	EnumVaultSecretEngineKv1
	EnumVaultSecretGeneric
)

func EnumVaultSecretEngineName() *wrapper.Enum[VaultSecretEngineName, string] {
	enum, _ := wrapper.NewEnum[VaultSecretEngineName](map[VaultSecretEngineName]string{
		EnumVaultSecretGeneric:   "generic",
		EnumVaultSecretEngineKv1: "kv-v1",
		EnumVaultSecretEngineKv2: "kv-v2",
	})
	return enum
}

type Secret interface {
	Render() SecretYaml
}

type Secrets map[string]Secret

func (s *Secrets) Render() *SecretsYaml {
	secretsYaml := SecretsYaml{}
	for k, v := range *s {
		secretsYaml[k] = v.Render()
	}
	return &secretsYaml
}

type SecretYaml interface {
	SecretYaml() SecretYaml
}
type SecretsYaml map[string]SecretYaml
type VaultSecretEngine struct {
	Name VaultSecretEngineName
	Path string
}

type VaultSecret struct {
	Vault VaultSecretConfig `yaml:"vault"`
	Token string            `yaml:"token,omitempty"`
}

func NewVaultSecret(engineName VaultSecretEngineName, enginePath string, path string, field string, token string) *VaultSecret {
	return &VaultSecret{
		Vault: VaultSecretConfig{
			Engine: VaultSecretEngine{
				Name: engineName,
				Path: enginePath,
			},
			Path:  path,
			Field: field,
		},
		Token: token,
	}
}

func NewVaultKv2Secret(enginePath string, path string, field string, token string) *VaultSecret {
	return NewVaultSecret(EnumVaultSecretEngineKv2, enginePath, path, field, token)
}

func (v *VaultSecret) Render() SecretYaml {
	return NewVaultSecretYaml(v)
}

type VaultSecretConfig struct {
	Engine VaultSecretEngine
	Path   string
	Field  string
}

type VaultSecretEngineYaml struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

func NewVaultSecretEngineYaml(engine *VaultSecretEngine) *VaultSecretEngineYaml {
	return &VaultSecretEngineYaml{
		Name: EnumVaultSecretEngineName().ValueSafe(engine.Name),
		Path: engine.Path,
	}
}

type VaultSecretConfigYaml struct {
	Engine VaultSecretEngineYaml `yaml:"engine"`
	Path   string                `yaml:"path"`
	Field  string                `yaml:"field"`
}

type VaultSecretYaml struct {
	Vault VaultSecretConfigYaml `yaml:"vault"`
	Token string                `yaml:"token,omitempty"`
}

func (v *VaultSecretYaml) SecretYaml() SecretYaml {
	return v
}

func NewVaultSecretYaml(secret *VaultSecret) *VaultSecretYaml {
	return &VaultSecretYaml{
		Vault: VaultSecretConfigYaml{
			Engine: *NewVaultSecretEngineYaml(&secret.Vault.Engine),
			Path:   secret.Vault.Path,
			Field:  secret.Vault.Field,
		},
		Token: secret.Token,
	}
}
