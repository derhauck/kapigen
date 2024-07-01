package types

import (
	"fmt"
	"runtime"
)

type DetailedError struct {
	Msg      string
	Name     string
	Filename string
	Line     int
}

const delimiter = "\n======================================================================"

func (c *DetailedError) Full() string {
	return fmt.Sprintf("%s\n=> FILE:\t%s:%d\n=> FUNCTION:\t%s\n=> MESSAGE:\t%s%s", delimiter, c.Filename, c.Line, c.Name, c.Msg, delimiter)
}
func (c *DetailedError) Error() string {
	return c.Msg
}
func ErrorHandler(msg string, caller int) error {

	pc, filename, line, _ := runtime.Caller(caller)
	return &DetailedError{
		Msg:      msg,
		Name:     runtime.FuncForPC(pc).Name(),
		Filename: filename,
		Line:     line,
	}

}
func DetailedErrorE(err error) error {
	return ErrorHandler(err.Error(), 2)
}
func DetailedErrorf(msg string, a ...any) error {
	return ErrorHandler(fmt.Sprintf(msg, a...), 2)
}
func NewMissingArgError(name string) error {
	return ErrorHandler(fmt.Sprintf("missing config arg '%s', required", name), 2)
}
func NewMissingArgsError(names ...string) error {
	msg := "missing config args "
	for _, n := range names {
		msg += fmt.Sprintf("'%s', ", n)
	}
	msg += "at least one required"
	return ErrorHandler(msg, 2)
}
