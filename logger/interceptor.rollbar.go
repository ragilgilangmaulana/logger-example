package logger

import (
	"github.com/lunixbochs/vtclean"
	"github.com/luthfikw/govalidator"
	rollbar "github.com/rollbar/rollbar-go"
	"logging-example-go/cservice"
	"logging-example-go/serror"
	"logging-example-go/utils/utarray"
	"logging-example-go/utils/utstring"
)

type (
	// RollbarOptions type
	RollbarOptions struct {
		Key     string     `json:"key" valid:"required"`
		Name    string     `json:"name" valid:"required"`
		Token   string     `json:"token" valid:"required"`
		Version string     `json:"version" valid:"required"`
		Level   ErrorLevel `json:"level" valid:"required"`
	}

	// Rollbar interceptor
	Rollbar interface {
		LogInterceptor
		IsEnabled() bool
		Enable()
		Disable()
	}

	rollbarInterceptorObj struct {
		Level   ErrorLevel
		Enabled bool
	}
)

// RollbarInterceptor to create rollbar interceptor
func RollbarInterceptor(opt RollbarOptions) (obj Rollbar, errx serror.SError) {
	if ok, err := govalidator.ValidateStruct(opt); !ok {
		errx = serror.NewFromErrorc(err, "Invalid rollbar options")
		return obj, errx
	}

	rollbar.SetToken(opt.Token)
	rollbar.SetEnvironment(utstring.Env(cservice.AppEnv, cservice.EnvDevelopment))
	rollbar.SetCodeVersion(opt.Version)
	rollbar.SetServerHost(opt.Key)

	obj = &rollbarInterceptorObj{
		Level:   opt.Level,
		Enabled: true,
	}
	return obj, errx
}

func (rollbarInterceptorObj) Translate(lvl ErrorLevel, obj interface{}) string {
	return DefaultTranslate(lvl, obj)
}

func (ox rollbarInterceptorObj) Process(lvl ErrorLevel, msg string) {
	DefaultProcess(lvl, msg)

	rlvl := map[ErrorLevel]string{
		ErrorLevelLog:      rollbar.DEBUG,
		ErrorLevelInfo:     rollbar.INFO,
		ErrorLevelWarning:  rollbar.WARN,
		ErrorLevelCritical: rollbar.CRIT,
	}

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
			return
		}

		if utarray.IsExist(lvl, allow) {
			rollbar.Log(rlvl[lvl], vtclean.Clean(msg, false))
		}
	}
}

func (ox rollbarInterceptorObj) IsEnabled() bool {
	return ox.Enabled
}

func (ox *rollbarInterceptorObj) Enable() {
	ox.Enabled = true
}

func (ox *rollbarInterceptorObj) Disable() {
	ox.Enabled = false
}
