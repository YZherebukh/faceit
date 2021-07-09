//go:generate mockgen -source ../logger/logger.go -destination ../logger/mock/mock_logger.go -exec_only false

package logger

import (
	"context"
	"fmt"

	cont "github.com/faceit/test/contextvalue"
)

type log interface {
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Logger struct {
	log log
}

// New creates a new logger
func New(l log) Logger {
	return Logger{log: l}
}

// Errorf logs with the Error severity.
// Arguments are handled in the manner of fmt.Printf.
// firs parameter is context. Logger will try go get processID coried by context
// so all logs for same request will have same processID
func (l Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.log.Errorf(fmt.Sprintf("processID:%s, %s", cont.ProcessID(ctx), format), v...)
}

// Infof logs with the Info severity.
// Arguments are handled in the manner of fmt.Printf.
// firs parameter is context. Logger will try go get processID coried by context
// so all logs for same request will have same processID
func (l Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.log.Infof(fmt.Sprintf("processID:%s, %s", cont.ProcessID(ctx), format), v...)
}

// Warningf logs with the Warning severity.
// Arguments are handled in the manner of fmt.Printf.
// firs parameter is context. Logger will try go get processID coried by context
// so all logs for same request will have same processID
func (l Logger) Warningf(ctx context.Context, format string, v ...interface{}) {
	l.log.Warningf(fmt.Sprintf("processID:%s, %s", cont.ProcessID(ctx), format), v...)
}

// Fatalf logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Printf.
// firs parameter is context. Logger will try go get processID coried by context
// so all logs for same request will have same processID
func (l Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.log.Fatalf(fmt.Sprintf("processID:%s, %s", cont.ProcessID(ctx), format), v...)
}
