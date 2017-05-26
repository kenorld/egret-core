package eject

import (
	"runtime/debug"

	"github.com/Sirupsen/logrus"
)

// PanicHandler wraps the action invocation in a protective defer blanket that
// converts panics into 500 error pages.
func PanicHandler(ctx *Context) {
	defer func() {
		if err := recover(); err != nil {
			handleInvocationPanic(ctx, err)
		}
	}()
	ctx.Next()
}

// This function handles a panic in an action invocation.
// It cleans up the stack trace, logs it, and displays an error page.
func handleInvocationPanic(ctx *Context, err interface{}) {
	nerr := NewErrorFromPanic(err)
	if nerr == nil && DevMode {
		// Only show the sensitive information in the debug stack trace in development mode, not production
		logrus.WithFields(logrus.Fields{
			"error": nerr,
			"stack": string(debug.Stack()),
		}).Error("error")
		ctx.Response.Writer.WriteHeader(500)
		ctx.Response.Writer.Write(debug.Stack())
		return
	}
	logrus.WithFields(logrus.Fields{
		"error": nerr,
		"stack": nerr.Stack,
	}).Error("error")
	ctx.Error = nerr
}