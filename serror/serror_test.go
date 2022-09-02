package serror

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicConstruction(t *testing.T) {
	errx := New("Test error message")
	_, _, line, _ := runtime.Caller(0)

	// expected
	var (
		exPath  = fmt.Sprintf("serror/serror_test.go:%d", (line - 1))
		exTitle = "Test error message"
		exMsg   = "Test error message"
	)

	if !assert.Contains(t, fmt.Sprintf("%s:%d", errx.File(), errx.Line()), exPath, "Error trace not matched") {
		t.FailNow()
	}

	if !assert.Equal(t, errx.Title(), exTitle, "Error title not matched") {
		t.FailNow()
	}

	if !assert.Equal(t, errx.Error(), exMsg, "Error message not matched") {
		t.FailNow()
	}
}

func TestBasicWithCommentConstruction(t *testing.T) {
	errx := Newc("Test error message", "test comments")
	_, _, line, _ := runtime.Caller(0)

	// expected
	var (
		exPath    = fmt.Sprintf("serror/serror_test.go:%d", (line - 1))
		exMsg     = "Test error message"
		exTitle   = "Test comments"
		exComment = "Test comments"
	)

	if !assert.Contains(t, fmt.Sprintf("%s:%d", errx.File(), errx.Line()), exPath, "Error trace not matched") {
		t.FailNow()
	}

	if !assert.Equal(t, errx.Title(), exTitle, "Error title not matched") {
		t.FailNow()
	}

	if !assert.Equal(t, errx.Comments(), exComment, "Error comments not matched") {
		t.FailNow()
	}

	if !assert.Equal(t, errx.Error(), exMsg, "Error message not matched") {
		t.FailNow()
	}
}
