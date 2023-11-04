package los

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kapigen.kateops.com/internal/logger"
	"net/http"
)

const LosHostName = "los.kateops.com"

type Client struct {
	client   *http.Client
	Hostname string
	Token    string
}

func NewClient(hostname string, token string) *Client {
	return &Client{
		client:   &http.Client{},
		Hostname: hostname,
		Token:    token,
	}
}

func (c *Client) GetLatestVersion(projectId string, path string) string {
	var defaultVersion = "0.0.0"
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://%s/v1/projects/%s/versions/latest?path=%s", c.Hostname, projectId, path),
		nil)
	req.Header.Set("AUTH", c.Token)
	resp, err := c.client.Do(req)

	if resp.StatusCode == 404 {
		return defaultVersion
	}

	if resp.StatusCode > 399 {
		logger.Error(resp.Status)
		var errResponse DefaultErrorResponse
		_ = ReadBody(resp.Body, &errResponse)
		logger.Error(errResponse.Message)
		return defaultVersion
	}

	var versionResponse ProjectVersionResponse
	err = ReadBody(resp.Body, &versionResponse)
	if err != nil {
		logger.ErrorE(err)
		return defaultVersion
	}
	return versionResponse.Version
}

func (c *Client) SetVersion(projectId string, version string, path string, tag string, latest bool) bool {
	var versionEndpoint = "https://%s/v1/projects/%s/versions"
	if latest {
		versionEndpoint = "https://%s/v1/projects/%s/versions/latest"
	}
	body, err := json.Marshal(projectSetVersion{
		Version: version,
		Path:    path,
		Tag:     tag,
	})
	if err != nil {
		logger.ErrorE(err)
		return false
	}
	req, err := http.NewRequest("POST",
		fmt.Sprintf(versionEndpoint, c.Hostname, projectId),
		bytes.NewReader(body),
	)
	req.Header.Set("AUTH", c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		logger.ErrorE(err)
		return false
	}
	if res.StatusCode < 200 || res.StatusCode > 399 {
		logger.Error(res.Status)
		var errResponse DefaultErrorResponse
		_ = ReadBody(res.Body, &errResponse)
		logger.Error(errResponse.Message)
		return false
	}
	return true
}
