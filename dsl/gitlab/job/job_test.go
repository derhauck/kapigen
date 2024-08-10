package job

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"gitlab.com/kateops/kapigen/dsl/gitlab/images"
	"gitlab.com/kateops/kapigen/dsl/gitlab/stages"
)

func TestNewCiJobYaml(t *testing.T) {
	t.Run("can create new CiJobYaml from CiJob", func(t *testing.T) {
		job := &CiJob{
			Script: NewScript(),
			Cache:  NewCache(),
			Image: Image{
				Name:       "image",
				PullPolicy: images.Always,
			},
			Stage:    stages.NewStage(),
			Services: Services{},
			Coverage: "coverage",
			Rules:    *DefaultMainBranchRules(),
			Secrets: Secrets{
				"TEST": &VaultSecret{
					VaultSecretConfig{
						Engine: VaultSecretEngine{
							Name: EnumVaultSecretEngineKv2,
							Path: "mount",
						},
						Path:  "path",
						Field: "field",
					},
					"token",
				},
			},
		}
		yaml, err := job.Render(nil, nil)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		snaps.MatchSnapshot(t, job, yaml)
	})
}
