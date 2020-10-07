package main

import (
	"context"
	"net/http"
	"net/http/httputil"
)

var ctx = context.Background()

// ABTestingProxy ...
type ABTestingProxy struct {
	Config            ABTestingProxyConfig
	DefaultProxy      httputil.ReverseProxy
	AdditionalProxies map[string]httputil.ReverseProxy
}

func (p *ABTestingProxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	proxy := p.DefaultProxy
	cookies := req.Cookies()
	for _, c := range cookies {
		if c.Name == p.Config.CookieName {
			if p, ok := p.AdditionalProxies[c.Value]; ok {
				proxy = p
			}
		}
	}
	proxy.ServeHTTP(wr, req)
}
