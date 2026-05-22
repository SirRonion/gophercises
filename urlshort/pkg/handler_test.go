package urlshort

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestYAMLHandlert(t *testing.T) {
	yamlData := []byte(`
- path: /google
  url: https://google.com
- path: /github
  url: https://github.com
`)

	fallbackCalled := false
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fallbackCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler, err := YAMLHandler(yamlData, fallback)
	if err != nil {
		t.Fatalf("Failed to build YAMLHandler: %v", err)
	}

	t.Run("Redirect Path", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/google", nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		if res.Code != http.StatusFound {
			t.Errorf("Expected status 302, got %d", res.Code)
		}

		location := res.Header().Get("Location")
		if location != "https://google.com" {
			t.Errorf("Expected location 'https://google.com', got '%s'", location)
		}
	})

	t.Run("Fallback Path", func(t *testing.T) {
		fallbackCalled = false
		req := httptest.NewRequest("GET", "/unkown-path", nil)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		if !fallbackCalled {
			t.Error("Expected fallback handler to be called, but it was not")
		}
		if res.Code != http.StatusOK {
			t.Errorf("Expected status 200 from fallback, got %d", res.Code)
		}
	})
}
