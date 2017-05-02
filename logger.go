package hooks

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Logger func(message string)

func NewLogger(w io.Writer) Logger {
	lastLog := time.Now()
	loggerOuput := log.New(w, "", log.LstdFlags)
	distanceDuration := 850 * time.Millisecond
	return func(logMessage string) {
		nowLog := time.Now()
		if nowLog.Before(lastLog.Add(distanceDuration)) {
			// don't use the log.Logger to print this message
			// if the last one was printed before some seconds.
			fmt.Println(logMessage) // fmt because we don't want the time, dev is dev so console.
			lastLog = nowLog
			return
		}
		// begin with new line in order to have the time once at the top
		// and the child logs below it.
		loggerOuput.Println("\u2192\n" + logMessage)
		lastLog = nowLog
	}
}

func DefaultLogger() Logger {
	return NewLogger(os.Stdout)
}
