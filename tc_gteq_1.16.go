//go:build go1.16
// +build go1.16

package tc

import (
	"io"
	"log"
)

func setDummyLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}
