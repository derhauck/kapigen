package los

import (
	"encoding/json"
	"io"
	"kapigen.kateops.com/internal/logger"
)

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
