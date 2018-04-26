package logrusmiddleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	h "github.com/bakins/test-helpers"
)

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func TestHandler(t *testing.T) {
	var buf bytes.Buffer

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "Hello World\n")
	})

	logger := logrus.New()
	logger.Level = logrus.InfoLevel
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = &buf

	l := Middleware{
		Name:   "example",
		Logger: logger,
	}

	lh := l.Handler(http.HandlerFunc(handler), "homepage")
	http.Handle("/", lh)

	lh.ServeHTTP(httptest.NewRecorder(), newRequest("GET", "/foo"))

	h.Assert(t, buf.Len() > 0, "buffer should not be empty")
	h.Assert(t, strings.Contains(buf.String(), `"component":"homepage"`), "buffer did not match expected result")
	var lr interface{}
	err := json.Unmarshal([]byte(buf.String()), &lr)
	h.Assert(t, err == nil, "could not decode log record")
	h.Assert(t, lr.(map[string]interface{})["status"].(float64) == 200, "should have status field set in log to match response code")
}
