package version

import (
	"errors"
	"fmt"
	"os"

	"kapigen.kateops.com/internal/logger"
)

const DefaultTagFile = ".version"

type FileReader struct {
	current map[string]string
}

func NewFileReader() *FileReader {
	return &FileReader{current: make(map[string]string)}
}
func (r *FileReader) read(path string) error {
	if path == "" {
		return errors.New("no path to read tag")
	}
	fileName := fmt.Sprintf("%s/%s", path, DefaultTagFile)
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	r.current[path] = string(file)
	return nil
}

func (r *FileReader) Get(path string) string {
	if r.current[path] == "" {
		if err := r.read(path); err != nil {
			logger.ErrorE(err)
			return ""
		}
	}
	return r.current[path]
}
