package log

import (
	"os"

	"github.com/go-kit/kit/log"
)

// logger is initialized by `func init` then used by `func Log` to print
// messages in a structured format.
var logger log.Logger

func init() {

	// Use a logfmt logger that writes to stderr.
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	// Add some default key/value data for all calls to Log.
	// User Caller(4) to expose the caller of this function.
	logger = log.WithPrefix(logger, "ts", log.DefaultTimestampUTC, "caller", log.Caller(4))
}

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
