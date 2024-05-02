package version

import (
	"os"
	"testing"
)

func TestFileReader_Get(t *testing.T) {
	reader := FileReader{
		current: make(map[string]string),
	}
	fullFileName := "./.version"
	t.Run("can read tag", func(t *testing.T) {
		err := os.WriteFile(fullFileName, []byte("1.0.0"), 0666)
		if err != nil {
			t.Errorf(err.Error())
		}

		reader.Get(".")
		if reader.current["."] != "1.0.0" {
			t.Errorf("can not read tag from file %s - %s, should be %s", ".", reader.current["."], "1.0.0")
		}
	})
	t.Cleanup(func() {
		err := os.Remove(fullFileName)
		if err != nil {
			t.Errorf(err.Error())
		}
	})
}
