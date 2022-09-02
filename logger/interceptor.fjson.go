package logger

import (
	"encoding/json"
	"fmt"
	"logging-example-go/serror"
	"strings"
	"time"

	"github.com/luthfikw/govalidator"
)

type (
	// FJSONOptions type
	FJSONOptions struct {
		Key      string `json:"key" valid:"required"`
		Name     string `json:"name" valid:"required"`
		Version  string `json:"version" valid:"required"`
		Printing bool   `json:"printing"`
	}

	// FJSON interceptor
	FJSON interface {
		LogInterceptor
		IsEnabled() bool
		Enable()
		Disable()
	}

	fjsonInterceptorObj struct {
		Key      string
		Name     string
		Version  string
		Printing bool
		Enabled  bool
	}
)

// FJSONInterceptor to create file json interceptor
func FJSONInterceptor(opt FJSONOptions) (obj FJSON, errx serror.SError) {
	if ok, err := govalidator.ValidateStruct(opt); !ok {
		errx = serror.NewFromErrorc(err, "Invalid fjson options")
		return obj, errx
	}

	obj = &fjsonInterceptorObj{
		Key:      opt.Key,
		Name:     opt.Name,
		Version:  opt.Version,
		Printing: opt.Printing,
		Enabled:  true,
	}
	return obj, errx
}

func (ox fjsonInterceptorObj) Translate(lvl ErrorLevel, obj interface{}) string {
	data := struct {
		Key       string      `json:"key"`
		Name      string      `json:"name"`
		Version   string      `json:"version"`
		Level     string      `json:"level"`
		Timestamp time.Time   `json:"timestamp"`
		Payload   interface{} `json:"payload"`
	}{
		Key:       ox.Key,
		Name:      ox.Name,
		Version:   ox.Version,
		Level:     strings.ToUpper(string(lvl)),
		Timestamp: time.Now(),
		Payload:   obj,
	}

	byt, err := json.Marshal(data)
	if err != nil {
		DefaultProcess(ErrorLevelCritical, fmt.Sprintf("Failed parsing data, details: %v", err))
		return ""
	}

	return string(byt)
}

func (ox fjsonInterceptorObj) Process(lvl ErrorLevel, msg string) {
	if ox.Printing {
		DefaultProcess(lvl, msg)
	}
}

func (ox fjsonInterceptorObj) IsEnabled() bool {
	return ox.Enabled
}

func (ox *fjsonInterceptorObj) Enable() {
	ox.Enabled = true
}

func (ox *fjsonInterceptorObj) Disable() {
	ox.Enabled = false
}

func (ox *fjsonInterceptorObj) StopPrinting() {
	ox.Printing = false
}

func (ox *fjsonInterceptorObj) StartPrinting() {
	ox.Printing = true
}
