package qrn

import (
	"fmt"
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type QueryLog struct {
	Query string
	Time  time.Duration
}

type Logger struct {
	Channel chan QueryLog
}

type ClosableDiscard struct{}

func (_ *ClosableDiscard) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (_ *ClosableDiscard) Close() error {
	return nil
}

func NewLogger(out io.WriteCloser) *Logger {
	ch := make(chan QueryLog)

	logger := &Logger{
		Channel: ch,
	}

	go func() {
		for ql := range ch {
			log, _ := jsoniter.MarshalToString(ql)
			fmt.Fprintln(out, log)
		}

		out.Close()
	}()

	return logger
}

func (logger *Logger) Log(query string, time time.Duration) {
	ql := QueryLog{
		Query: query,
		Time:  time,
	}

	logger.Channel <- ql
}

func (logger *Logger) Close() {
	close(logger.Channel)
}
