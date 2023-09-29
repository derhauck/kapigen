package logger

import (
	"fmt"
	"os"
)

type level struct {
	Info  string
	Debug string
	Error string
}

var Level = &level{
	Info:  "INFO",
	Error: "ERROR",
	Debug: "DEBUG",
}

func sortLevel(levelA string, levelB string) bool {
	switch levelA {

	case Level.Error:
		if levelB == Level.Error {
			return true
		} else {
			return false
		}
	case Level.Info:
		if levelB == Level.Error || levelB == Level.Info {
			return true
		} else {
			return false
		}
	case Level.Debug:
		return true
	default:
		return false
	}
}

func escalate(level string) bool {
	currentLevel := os.Getenv("LOGGER_LEVEL")
	return sortLevel(currentLevel, level)
}

func log(level string, msg ...string) {
	if escalate(level) {
		for i := 0; i < len(msg); i++ {
			fmt.Println(fmt.Sprintf("%s:\t%s", level, msg[i]))
		}
	}

}

func logAny(level string, msg any) {
	if escalate(level) {
		fmt.Println(fmt.Sprintf("%s: %s", level, fmt.Sprint(msg)))
	}
}

func Info(msg ...string) {
	log(Level.Info, msg...)
}

func Debug(msg ...string) {
	log(Level.Debug, msg...)
}

func DebugAny(msg any) {

	logAny(Level.Debug, msg)
}
func Error(msg ...string) {
	log(Level.Error, msg...)
}

func ErrorE(err error) {
	log(Level.Error, err.Error())
}
