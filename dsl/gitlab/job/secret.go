package job

import (
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

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

type VaultSecretEngine struct {
	Name VaultSecretEngineName
	Path string
}

type VaultSecret struct {
	vault VaultSecretConfig `yaml:"vault"`
}

type VaultSecretConfig struct {
	Engine VaultSecretEngine
	Path   string
	Field  string
	Token  string
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

type VaultSecretYaml struct {
	Engine VaultSecretEngineYaml `yaml:"engine"`
	Path   string                `yaml:"path"`
	Field  string                `yaml:"field"`
	Token  string                `yaml:"token,omitempty"`
}

func NewVaultSecretYaml(secret *VaultSecret) *VaultSecretYaml {
	return &VaultSecretYaml{
		Engine: *NewVaultSecretEngineYaml(&secret.Engine),
		Path:   secret.Path,
		Field:  secret.Field,
		Token:  secret.Token,
	}
}
