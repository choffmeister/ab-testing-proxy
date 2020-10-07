package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
	"testing"
)

type testIdentifier struct{}

func (t testIdentifier) IdentifyRequest(req *http.Request) (*string, error) {
	id := req.URL.Query().Get("id")
	if id == "" {
		return nil, fmt.Errorf("anonymous")
	}
	return &id, nil
}

func TestABTestingProxy(t *testing.T) {
	defaultProxy, _ := CreateProxyTarget("http://localhost:10000")
	otherProxy, _ := CreateProxyTarget("http://localhost:10001")
	proxy := ABTestingProxy{
		Config: ABTestingProxyConfig{
			CookieName: "test",
			DefaultTarget: ABTestingProxyConfigTarget{
				Id:  "default",
				Url: "http://localhost:10000",
			},
			AdditionalTargets: []ABTestingProxyConfigTarget{
				ABTestingProxyConfigTarget{
					Id:  "other",
					Url: "http://localhost:10001",
				},
			},
		},
		DefaultProxy: *defaultProxy,
		AdditionalProxies: map[string]httputil.ReverseProxy{
			"other": *otherProxy,
		},
	}

	test := func(path string, cookie string, expectedResponse string) {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		if cookie != "" {
			req.Header.Add("Cookie", fmt.Sprintf("test=%s", cookie))
		}

		w := httptest.NewRecorder()
		proxy.ServeHTTP(w, req)
		b := strings.Trim(w.Body.String(), " \n")
		if b != expectedResponse {
			t.Errorf("expected response %s, got %s", expectedResponse, b)
		}
	}

	t.Run("without cookie", func(t *testing.T) {
		test("/", "", "backend-1")
	})

	t.Run("with invalid cookie", func(t *testing.T) {
		test("/", "unknown", "backend-1")
	})

	t.Run("with cookie", func(t *testing.T) {
		test("/", "default", "backend-1")
		test("/", "other", "backend-2")
	})
}
