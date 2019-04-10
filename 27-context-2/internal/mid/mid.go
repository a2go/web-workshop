package mid

import (
	"log"
)

// Middleware holds the required state for all web.Middleware functions in this
// package. Its methods are defined in separate files.
type Middleware struct {
	Log *log.Logger
}
