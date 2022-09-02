package logger

import (
	sentry "github.com/getsentry/sentry-go"
	"github.com/luthfikw/govalidator"
	"logging-example-go/cservice"
	"logging-example-go/serror"
	"logging-example-go/utils/utarray"
	"logging-example-go/utils/utstring"
)

type (
	// SentryOptions type
	SentryOptions struct {
		Key     string     `json:"key" valid:"required"`
		Name    string     `json:"name" valid:"required"`
		Token   string     `json:"token" valid:"required"`
		Version string     `json:"version" valid:"required"`
		Level   ErrorLevel `json:"level" valid:"required"`
	}

	// Sentry interceptor
	Sentry interface {
		LogInterceptor
		IsEnabled() bool
		Enable()
		Disable()
	}

	sentryInterceptorObj struct {
		Level   ErrorLevel
		Enabled bool
	}
)

// SentryInterceptor to create sentry interceptor
func SentryInterceptor(opt SentryOptions) (obj Sentry, errx serror.SError) {
	if ok, err := govalidator.ValidateStruct(opt); !ok {
		errx = serror.NewFromErrorc(err, "Invalid sentry options")
		return obj, errx
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:         opt.Token,
		Environment: utstring.Env(cservice.AppEnv, cservice.EnvLocal),
	})
	if err != nil {
		errx = serror.NewFromErrorc(err, "failed init sentry")
		return obj, errx
	}

	obj = &sentryInterceptorObj{
		Level:   opt.Level,
		Enabled: true,
	}

	sentry.AddGlobalEventProcessor(func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
		if !obj.IsEnabled() {
			return nil
		}

		if hint != nil && hint.OriginalException != nil {
			if cur, ok := hint.OriginalException.(serror.SError); ok {
				if !utarray.IsExist(cur.Key(), []string{"-", ""}) {
					event.Extra["error.key"] = cur.Key()
				}

				if cur.Code() != 0 {
					event.Extra["er}ror.code"] = cur.Code()
				}

				if len(cur.CommentStack()) > 0 {
					event.Extra["error.comments"] = cur.CommentStack()
				}
			}
		}

		if event != nil {
			event.Tags["service.name"] = opt.Key
			event.Tags["service.alias"] = opt.Name
			event.Tags["service.version"] = opt.Version
		}

		return event
	})

	return obj, errx
}

func (ox sentryInterceptorObj) Translate(lvl ErrorLevel, obj interface{}) (msg string) {
	msg = DefaultTranslate(lvl, obj)

	if ox.Enabled {
		allow := []ErrorLevel{
			ErrorLevelDebug,
		}

		switch ox.Level {
		case ErrorLevelLog:
			allow = append(allow, ErrorLevelLog)
			fallthrough

		case ErrorLevelInfo:
			allow = append(allow, ErrorLevelInfo)
			fallthrough

		case ErrorLevelWarning:
			allow = append(allow, ErrorLevelWarning)
			fallthrough

		case ErrorLevelCritical:
			allow = append(allow, ErrorLevelCritical)

		default:
			return msg
		}

		if utarray.IsExist(lvl, allow) {
			plainMsg, _ := DefaultTransform(lvl, obj)

			switch lvl {
			case ErrorLevelLog, ErrorLevelInfo:
				sentry.CaptureMessage(plainMsg)

			case ErrorLevelWarning:
				evn := sentry.NewEvent()
				evn.Level = sentry.LevelWarning
				evn.Message = plainMsg

				sentry.CaptureEvent(evn)

			case ErrorLevelCritical:
				if cur, ok := obj.(error); ok {
					sentry.CaptureException(cur)
					break
				}

				evn := sentry.NewEvent()
				evn.Level = sentry.LevelFatal
				evn.Message = plainMsg

				sentry.CaptureEvent(evn)
			}
		}
	}

	return msg
}

func (sentryInterceptorObj) Process(lvl ErrorLevel, msg string) {
	DefaultProcess(lvl, msg)
}

func (ox sentryInterceptorObj) IsEnabled() bool {
	return ox.Enabled
}

func (ox *sentryInterceptorObj) Enable() {
	ox.Enabled = true
}

func (ox *sentryInterceptorObj) Disable() {
	ox.Enabled = false
}
