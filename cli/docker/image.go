package docker

type Image int

const (
	Kapigen_Latest Image = iota
	Alpine_3_18
	Terraform_Base
)

func (c Image) Image() string {
	return []string{
		"kapigen",
		"alpine:3.18",
		"image hub.kateops.com/base/terraform:latest",
	}[c]
}
