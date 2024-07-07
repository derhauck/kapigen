package logger

import (
	"fmt"

	"gitlab.com/kateops/kapigen/dsl/logger/level"
)

func escalate(level level.Level) bool {
	return level.IsActive()
}

func log(level level.Level, msg ...string) {
	if escalate(level) {
		for i := 0; i < len(msg); i++ {
			fmt.Printf("%s:\t%s\n", level, msg[i])
		}
	}

}

func logAny(level level.Level, msg any) {
	if escalate(level) {
		fmt.Printf("%s:\t%s\n", level, fmt.Sprint(msg))
	}
}

func Info(msg ...string) {
	log(level.Info, msg...)
}

func Debug(msg ...string) {
	log(level.Debug, msg...)
}

func DebugAny(msg any) {

	logAny(level.Debug, msg)
}
func Error(msg ...string) {
	log(level.Error, msg...)
}

func ErrorE(err error) {
	log(level.Error, err.Error())
}
