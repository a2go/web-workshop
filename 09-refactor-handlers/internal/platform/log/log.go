/*
Package log provides an implementation of a log function with no knobs that
this project uses for all internal logging purposes.
*/
package log

import (
	"os"

	"github.com/go-kit/kit/log"
)

var logger = func() log.Logger {
	var lg = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	// User Caller(4) to expose the caller of this function.
	lg = log.WithPrefix(lg, "ts", log.DefaultTimestampUTC, "caller", log.Caller(4))

	return lg
}()

// Log a human readible message with a variadic sequence of alternating
// key-value pairs. Output format is `logfmt` (See: https://brandur.org/logfmt).
// Output is sent to stderr.
func Log(msg string, keyvals ...interface{}) error {
	return log.With(logger, "msg", msg).Log(keyvals...)
}

/* Non-structured alternative
var logger = log.New(os.Stderr, "", log.LstdFlags)

func Log(format string, v ...interface{}) {
	logger.Printf(format, v...)
}
*/
