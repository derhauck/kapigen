package docker

type Image int

const (
	Kapigen_Latest Image = iota
	Alpine_3_18
	Terraform_Base
	Buildkit
	BUILDKITD
)

var values = map[Image]string{
	Kapigen_Latest: "kapigen",
	Alpine_3_18:    "alpine:3.18",
	Terraform_Base: "hub.kateops.com/base/terraform:latest",
	Buildkit:       "moby/buildkit:master",
	BUILDKITD:      "moby/buildkit:master-rootless",
}

func (c Image) Image() string {
	if value, ok := values[c]; ok {
		return value
	}
	return values[Alpine_3_18]
}
