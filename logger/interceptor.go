package logger

import (
	"fmt"
	"logging-example-go/serror"
	"logging-example-go/utils/utstring"
	"logging-example-go/utils/uttime"
	"os"
	"runtime"
	"time"
)

type (
	// LogInterceptor interface
	LogInterceptor interface {
		Translate(lvl ErrorLevel, obj interface{}) string
		Process(lvl ErrorLevel, msg string)
	}

	defaultInterceptorObj struct{}
)

func (defaultInterceptorObj) Translate(lvl ErrorLevel, obj interface{}) string {
	return DefaultTranslate(lvl, obj)
}

func (defaultInterceptorObj) Process(lvl ErrorLevel, msg string) {
	DefaultProcess(lvl, msg)
}

// DefaultInterceptor for default interceptor
func DefaultInterceptor() LogInterceptor {
	return defaultInterceptorObj{}
}

// DefaultTranslate for default translation
func DefaultTranslate(lvl ErrorLevel, obj interface{}) string {
	m, m2 := DefaultTransform(lvl, obj)

	if !isLocal() {
		m2 = m
	}

	// formating
	cur := time.Now()

	type lvll struct {
		Label string
		Color utstring.Color
	}

	lbl := "?"
	ls := map[ErrorLevel]lvll{
		ErrorLevelInfo:     lvll{"INFO", utstring.LIGHT_BLUE},
		ErrorLevelLog:      lvll{"LOG", utstring.LIGHT_GRAY},
		ErrorLevelWarning:  lvll{"WARN", utstring.LIGHT_YELLOW},
		ErrorLevelCritical: lvll{"ERR", utstring.RED},
	}
	if cur, ok := ls[lvl]; ok {
		lbl = cur.Label

		if isLocal() {
			lbl = utstring.ApplyForeColor(lbl, cur.Color)
		}
	}

	return fmt.Sprintf("[%s] %s: %s", uttime.Format(uttime.DefaultDateTimeFormat, cur), lbl, m2)
}

// DefaultTransform for default transforming
func DefaultTransform(lvl ErrorLevel, obj interface{}) (plainMsg string, colorMsg string) {
	plainMsg = fmt.Sprintf("%v", obj)
	colorMsg = plainMsg

	switch lvl {
	case ErrorLevelCritical, ErrorLevelWarning:
		switch vx := obj.(type) {
		case serror.SError:
			plainMsg = vx.String()
			colorMsg = vx.ColoredString()

		case error:
			pc, fn, line, _ := runtime.Caller(4)
			plainMsg = fmt.Sprintf(serror.StandardFormat(), runtime.FuncForPC(pc).Name(), fn, line, plainMsg)
			colorMsg = fmt.Sprintf(serror.StandardColorFormat(), runtime.FuncForPC(pc).Name(), fn, line, colorMsg)
		}
	}

	return plainMsg, colorMsg
}

// DefaultProcess for default processing
func DefaultProcess(lvl ErrorLevel, msg string) {
	if msg == "" {
		return
	}

	switch lvl {
	case ErrorLevelCritical, ErrorLevelWarning:
		DefaultStderr(msg)

	default:
		DefaultStdout(msg)
	}
}

// DefaultStdout for default stdout print
func DefaultStdout(msg string) {
	fmt.Fprintln(os.Stdout, msg)
}

// DefaultStdout for default stderr print
func DefaultStderr(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}
