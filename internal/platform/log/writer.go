package log

import (
	"bufio"
	"io"
	"runtime"

	"github.com/InVisionApp/go-logger"
)

// NewWriter creates a new writer from a Logger.
func NewWriter(logger log.Logger, level Level) *io.PipeWriter {
	reader, writer := io.Pipe()

	var printFunc func(args ...interface{})

	switch level {
	case DebugLevel:
		printFunc = logger.Debug
	case InfoLevel:
		printFunc = logger.Info
	case WarnLevel:
		printFunc = logger.Warn
	case ErrorLevel:
		printFunc = logger.Error
	default:
		printFunc = logger.Info
	}

	go writerScanner(logger, reader, printFunc)
	runtime.SetFinalizer(writer, writerFinalizer)

	return writer
}

// nolint: interfacer
func writerScanner(logger log.Logger, reader *io.PipeReader, printFunc func(args ...interface{})) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		printFunc(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		logger.Errorf("Error while reading from Writer: %s", err)
	}
	reader.Close()
}

func writerFinalizer(writer *io.PipeWriter) {
	writer.Close()
}
