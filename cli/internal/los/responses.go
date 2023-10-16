package los

type ProjectVersionResponse struct {
	Created string `json:"created"`
	Path    string `json:"path"`
	Tag     string `json:"tag"`
	Version string `json:"version"`
}

type DefaultErrorResponse struct {
	Message string `json:"message"`
}
