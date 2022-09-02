package logger

import (
	"fmt"
	"logging-example-go/cservice"
	"logging-example-go/serror"
	"logging-example-go/utils/utstring"
	"os"
	"strings"
	"syscall"
)

var interceptor LogInterceptor = DefaultInterceptor()

// SetInterceptor to set log interceptor
func SetInterceptor(inc LogInterceptor) {
	if inc == nil {
		panic("Null interceptor")
	}

	interceptor = inc
}

// Info to logging info level
func Info(msg interface{}) {
	lvl := ErrorLevelInfo

	m := interceptor.Translate(lvl, msg)
	interceptor.Process(lvl, m)
}

// Infof to logging info level with function
func Infof(msg string, args ...interface{}) {
	Info(fmt.Sprintf(msg, args...))
}

// Log to logging log level
func Log(msg interface{}) {
	lvl := ErrorLevelLog

	m := interceptor.Translate(lvl, msg)
	interceptor.Process(lvl, m)
}

// Logf to logging log level with function
func Logf(msg string, args ...interface{}) {
	Log(fmt.Sprintf(msg, args...))
}

// Warn to logging warning level
func Warn(msg interface{}) {
	lvl := ErrorLevelWarning

	m := interceptor.Translate(lvl, msg)
	interceptor.Process(lvl, m)
}

// Warnf to logging warning level with function
func Warnf(msg string, args ...interface{}) {
	Warn(fmt.Sprintf(msg, args...))
}

// Err to logging error level
func Err(msg interface{}) {
	lvl := ErrorLevelCritical

	m := interceptor.Translate(lvl, msg)
	interceptor.Process(lvl, m)
}

// Errf to logging error level with function
func Errf(msg string, args ...interface{}) {
	Err(serror.Newsf(1, msg, args...))
}

// Panic to logging error then exit
func Panic(msg interface{}) {
	Err(castToSError(msg, 1))
	exit()
}

// private

func isLocal() bool {
	return strings.ToLower(utstring.Env(cservice.AppEnv, cservice.EnvLocal)) == cservice.EnvLocal
}

func exit() {
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if err != nil {
		os.Exit(1)
	}
}

func castToSError(obj interface{}, skip int) serror.SError {
	var errx serror.SError

	if cur, ok := obj.(serror.SError); ok {
		errx = cur

	} else if cur, ok := obj.(error); ok {
		errx = serror.NewFromErrors(skip+1, cur)

	} else {
		errx = serror.Newsf(skip+1, "%+v", obj)
	}

	return errx
}
