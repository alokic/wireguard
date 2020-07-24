package testhelpers

import (
	"bytes"
	"fmt"

	"github.com/go-kit/kit/log"
)

type TestLogger struct {
	print bool
}

func (logger TestLogger) Log(args ...interface{}) error {
	if logger.print {
		fmt.Println(args...)
	}
	return nil
}

// If you want to print the logs pass print=true
// Eg. testhelpers.FakeLogger(true)
func FakeLogger(print bool) log.Logger {
	return TestLogger{print}
}

func LoggerWithWriter() (log.Logger, *bytes.Buffer) {
	writer := &bytes.Buffer{}
	return log.NewLogfmtLogger(writer), writer
}
