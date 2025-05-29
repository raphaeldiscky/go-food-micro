// Package testutils provides a test utils.
package testutils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"emperror.dev/errors"

	echo "github.com/labstack/echo/v4"
)

// SkipCI skips the test if the CI environment is set.
func SkipCI(t *testing.T) {
	t.Helper()

	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")

		return
	}
}

// WaitUntilConditionMet waits until the condition is met.
func WaitUntilConditionMet(conditionToMet func() bool, timeout ...time.Duration) error {
	timeOutTime := 20 * time.Second
	if len(timeout) >= 0 && timeout != nil {
		timeOutTime = timeout[0]
	}

	startTime := time.Now()
	timeOutExpired := false
	meet := conditionToMet()
	for !meet {
		if timeOutExpired {
			return errors.New("Condition not met for the test, timeout exceeded")
		}
		time.Sleep(time.Second * 2)
		meet = conditionToMet()
		timeOutExpired = time.Since(startTime) > timeOutTime
	}

	return nil
}

// HTTPRecorder records the http request and response.
func HTTPRecorder(
	t *testing.T,
	e *echo.Echo,
	req *http.Request,
	f func(w *httptest.ResponseRecorder) bool,
) {
	t.Helper()

	w := httptest.NewRecorder()
	e.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}
