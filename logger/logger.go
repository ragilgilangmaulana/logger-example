package logger

import (
	"errors"
	"fmt"
	"logging-example-go/serror"
	"logging-example-go/utils/utpath"
	"logging-example-go/utils/utstring"
	"logging-example-go/utils/uttime"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type (
	// ErrorLevel type
	ErrorLevel string
)

// Mode type
type Mode int

const (
	// ModeDaily mode
	ModeDaily Mode = 1 + iota

	// ModeMonthly mode
	ModeMonthly

	// ModeYearly mode
	ModeYearly

	// ModePermanent mode
	ModePermanent
)

const (
	// ErrorLevelDebug level
	ErrorLevelDebug ErrorLevel = "debug"

	// ErrorLevelLog level
	ErrorLevelLog ErrorLevel = "log"

	// ErrorLevelInfo level
	ErrorLevelInfo ErrorLevel = "info"

	// ErrorLevelCritical level
	ErrorLevelCritical ErrorLevel = "critical"

	// ErrorLevelWarning level
	ErrorLevelWarning ErrorLevel = "warn"
)

type (
	// Options struct
	Options struct {
		Mode        Mode
		Path        string
		Writing     bool
		FileFormat  string
		Interceptor LogInterceptor
	}

	logger struct {
		sync.Mutex
		Path       string
		Writing    bool
		Mode       Mode
		FileFormat string

		_ready       bool
		_file        string
		_name        string
		_queues      []string
		_interceptor LogInterceptor
		_stream      *os.File
	}

	// Logger service
	Logger interface {
		Startup() error
		Info(msg interface{})
		Infof(msg string, args ...interface{})
		Log(msg interface{})
		Logf(msg string, args ...interface{})
		Warn(msg interface{})
		Warnf(msg string, args ...interface{})
		Err(msg interface{})
		Errf(msg string, args ...interface{})
		Panic(msg interface{})
		IsWriting() bool
		StartWriting()
		StopWriting()
	}
)

// public

// Construct function
func Construct(opt Options) Logger {
	log := &logger{
		Mode:       opt.Mode,
		Path:       opt.Path,
		Writing:    opt.Writing,
		FileFormat: utstring.Chains(opt.FileFormat, "log-%v.log"),

		_ready:       false,
		_interceptor: opt.Interceptor,
	}
	return log
}

// Startup to starting up
func (ox *logger) Startup() error {
	var err error
	if ox.Writing {
		cur := time.Now()
		fv := map[string]string{
			"d": cur.Format("02"),
			"m": cur.Format("01"),
			"y": cur.Format("2006"),
			"h": cur.Format("15"),
			"i": cur.Format("04"),
			"s": cur.Format("05"),
			"v": "",
		}

		formt := ox.FileFormat

		switch ox.Mode {
		case ModeDaily:
			fv["v"] = cur.Format("20060102")

		case ModeMonthly:
			fv["v"] = cur.Format("200601fprintln")

		case ModeYearly:
			fv["v"] = cur.Format("2006")

		case ModePermanent:
			fv["v"] = ""
		}

		for k, v := range fv {
			formt = strings.ReplaceAll(formt, "%"+k, v)
		}

		if ox._name != formt {
			ox._name = formt
			ox._file = filepath.Join(ox.Path, ox._name)

			if !utpath.IsExists(ox.Path) {
				err = os.MkdirAll(ox.Path, os.ModePerm)
				if err != nil {
					return err
				}
			}

			err = ox.open()
			if err != nil {
				return err
			}
		}
	}

	if !ox._ready {
		go func() {
			for {
				time.Sleep(3 * time.Second)
				ox.flush()
			}
		}()
	}

	ox._ready = true
	return err
}

// Info to logging info level
func (ox *logger) Info(msg interface{}) {
	lvl := ErrorLevelInfo

	m := ox._interceptor.Translate(lvl, msg)
	ox._interceptor.Process(lvl, m)
	_ = ox.write(m)
}

// Infof to logging info level with function
func (ox *logger) Infof(msg string, args ...interface{}) {
	ox.Info(fmt.Sprintf(msg, args...))
}

// Log to logging log level
func (ox *logger) Log(msg interface{}) {
	lvl := ErrorLevelLog

	m := ox._interceptor.Translate(lvl, msg)
	ox._interceptor.Process(lvl, m)
	_ = ox.write(m)
}

// Logf to logging log level with function
func (ox *logger) Logf(msg string, args ...interface{}) {
	ox.Log(fmt.Sprintf(msg, args...))
}

// Warn to logging warning level
func (ox *logger) Warn(msg interface{}) {
	lvl := ErrorLevelWarning

	m := ox._interceptor.Translate(lvl, msg)
	ox._interceptor.Process(lvl, m)
	_ = ox.write(m)
}

// Warnf to logging warning level with function
func (ox *logger) Warnf(msg string, args ...interface{}) {
	ox.Warn(fmt.Sprintf(msg, args...))
}

// Err to logging error level
func (ox *logger) Err(msg interface{}) {
	lvl := ErrorLevelCritical

	m := ox._interceptor.Translate(lvl, msg)
	ox._interceptor.Process(lvl, m)
	_ = ox.write(m)
}

// Errf to logging error level with function
func (ox *logger) Errf(msg string, args ...interface{}) {
	ox.Err(serror.Newsf(1, msg, args...))
}

// Panic to logging error level then exit the app
func (ox *logger) Panic(msg interface{}) {
	ox.Err(castToSError(msg, 1))
	exit()
}

// IsReady to get ready flag
func (ox *logger) IsReady() bool {
	return ox._ready
}

// IsWriting to get writing flag
func (ox *logger) IsWriting() bool {
	return ox.Writing
}

// StopWriting to stop writing
func (ox *logger) StopWriting() {
	ox.Writing = false
}

// StartWriting to start writing
func (ox *logger) StartWriting() {
	ox.Writing = true
}

// private

func (ox *logger) open() error {
	if !ox.Writing {
		return nil
	}

	var err error

	ox.Lock()
	defer ox.Unlock()

	if ox._stream != nil {
		_ = ox._stream.Close()
	}

	ox._stream, err = os.OpenFile(ox._file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	return err
}

func (ox *logger) write(m string) error {
	if !ox.Writing {
		return nil
	}

	if !ox._ready {
		return errors.New("Logger not yet ready")
	}

	if m == "" {
		return nil
	}

	ox.Lock()
	ox._queues = append(ox._queues, m)
	ox.Unlock()

	return nil
}

func (ox *logger) flush() error {
	if !ox.Writing {
		return nil
	}

	err := ox.Startup()
	if err != nil {
		return err
	}

	ox.Lock()
	lists := ox._queues
	ox._queues = []string{}
	ox.Unlock()

	defer func() {
		if err != nil {
			ox.Lock()
			ox._queues = append(lists, ox._queues...)
			ox.Unlock()
		}
	}()

	if len(lists) > 0 {
		for _, v := range lists {
			_, err = ox._stream.WriteString(fmt.Sprintf("%s\n", v))
			if err != nil {
				ox.printf("Failed to writing, details: %+v", err)

				errs := ox.open()
				if errs != nil {
					ox.printf("Failed to re-open file %s, details: %+v", ox._file, errs)
				}
				return err
			}
		}

		err = ox._stream.Sync()
		if err != nil {
			ox.printf("Failed to flushing stream, details: %+v", err)
			return err
		}
	}

	return err
}

func (ox *logger) printf(msg string, opts ...interface{}) {
	fmt.Printf("[%s] ERR: %s\n", uttime.Format(uttime.DefaultDateTimeFormat, time.Now()), fmt.Sprintf(msg, opts...))
}
