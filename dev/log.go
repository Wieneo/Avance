package dev

import (
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

//DebugLogging should be set to true if app runs in debug mode
var (
	DeubgLogging = false
)

//LogInfo logs with the info prefix
func LogInfo(Message ...interface{}) {
	fmt.Printf("[%s] [I] %v\n", dateForLog(), Message)
}

//LogWarn logs with the info prefix
func LogWarn(Message ...interface{}) {
	fmt.Printf("[%s] [W] %v\n", dateForLog(), Message)
}

//LogError logs with the info prefix
func LogError(Error error, Message ...interface{}) {
	fmt.Printf("[%s] [ERROR] %v\n", dateForLog(), Message)
	sentry.CaptureException(Error)
	sentry.Flush(time.Second * 5)
}

//LogErrorNoSentry logs an error without informing sentry
func LogErrorNoSentry(Error error, Message ...interface{}) {
	fmt.Printf("[%s] [ERROR] %v\n", dateForLog(), Message)
}

//LogFatal logs with the info prefix
func LogFatal(Error error, Message ...interface{}) {
	fmt.Printf("[%s] [!FATAL!] %v\n", dateForLog(), Message)
	sentry.CaptureException(Error)
	sentry.Flush(time.Second * 5)
	os.Exit(255)
}

//LogDebug logs with the info prefix
func LogDebug(Message ...interface{}) {
	if DeubgLogging {
		fmt.Printf("[%s] [D] %v\n", dateForLog(), Message)
	}
}

func dateForLog() string {
	return fmt.Sprint(time.Now().Year(), "-", time.Now().Month(), "-", time.Now().Day(), " ", time.Now().Hour(), ":", time.Now().Minute(), ":", time.Now().Second())
}
