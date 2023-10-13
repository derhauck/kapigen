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

	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://los.kateops.com/v1/projects/%s/versions/latest?path=%s", project, path),
		nil)
	client := &http.Client{}
	req.Header.Set("AUTH", environment.LOS_AUTH.Get())
	resp, err := client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.ErrorE(err)
		}
	}(resp.Body)

	if resp.StatusCode == 404 {
		return "0.0.0"
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorE(err)
		return "0.0.0"
	}
	var versionResponse ProjectVersion
	err = json.Unmarshal(body, &versionResponse)
	if err != nil {
		logger.ErrorE(err)
	}
	return versionResponse.Version
}
