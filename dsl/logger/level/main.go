package level

import "os"

func GetCurrentLevel() Level {
	currentLevel := os.Getenv("LOGGER_LEVEL")
	if currentLevel == "" {
		return 0
	}
	level, err := FromString(currentLevel)
	if err != nil {
		return 0
	}
	return level
}
