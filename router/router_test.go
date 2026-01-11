package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicRouter(t *testing.T) {
	r := New()

	r.Handle("GET", "/hello", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("hello world!"))
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Errorf("expected 200, got %d", rec.Code)
	}

	if rec.Body.String() != "hello world!" {
		t.Errorf("expected hello world!, got %q", rec.Body.String())
	}
}
