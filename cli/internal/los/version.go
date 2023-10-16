package los

import (
	"encoding/json"
	"fmt"
	"io"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/logger"
	"net/http"
)

func GetVersion(project string, path string) string {
	var defaultVersion = "0.0.0"
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://los.kateops.com/v1/projects/%s/versions/latest?path=%s", project, path),
		nil)
	client := &http.Client{}
	req.Header.Set("AUTH", environment.LOS_AUTH_TOKEN.Get())
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.ErrorE(err)
		}
	}(resp.Body)

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorE(err)
		return defaultVersion
	}
	var versionResponse ProjectVersionResponse
	err = json.Unmarshal(body, &versionResponse)
	if err != nil {
		logger.ErrorE(err)
	}
	return versionResponse.Version
}

func ReadBody(body io.ReadCloser, v any) error {
	if body == nil {
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.ErrorE(err)
		}
	}(body)
	content, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, v)
}
