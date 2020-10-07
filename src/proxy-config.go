package main

import (
	utiljson "encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"net/http/httputil"
	"net/url"
)

// ABTestingProxyConfigTarget ...
type ABTestingProxyConfigTarget struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

// ABTestingProxyConfig ...
type ABTestingProxyConfig struct {
	CookieName        string                       `json:"cookieName"`
	DefaultTarget     ABTestingProxyConfigTarget   `json:"defaultTarget"`
	AdditionalTargets []ABTestingProxyConfigTarget `json:"additionalTargets"`
}

// LoadABTestingProxyConfig ...
func LoadABTestingProxyConfig(yaml []byte) (*ABTestingProxyConfig, *httputil.ReverseProxy, *map[string]httputil.ReverseProxy, error) {
	config := &ABTestingProxyConfig{}

	json, err := utilyaml.ToJSON(yaml)
	if err != nil {
		return nil, nil, nil, err
	}

	if err = utiljson.Unmarshal(json, config); err != nil {
		return nil, nil, nil, err
	}

	if config.CookieName == "" {
		return nil, nil, nil, fmt.Errorf("cookie name must not be empty")
	}

	defaultProxy, err := CreateProxyTarget(config.DefaultTarget.Url)
	if err != nil {
		return nil, nil, nil, err
	}

	additionalProxies := map[string]httputil.ReverseProxy{}
	for _, t := range config.AdditionalTargets {
		proxy, err := CreateProxyTarget(t.Url)
		if err != nil {
			return nil, nil, nil, err
		}
		additionalProxies[t.Id] = *proxy
	}

	return config, defaultProxy, &additionalProxies, nil
}

// Equal ...
func (c1 ABTestingProxyConfig) Equal(c2 ABTestingProxyConfig) bool {
	return true &&
		cmp.Equal(c1.CookieName, c2.CookieName) &&
		cmp.Equal(c1.DefaultTarget, c2.DefaultTarget) &&
		cmp.Equal(c1.AdditionalTargets, c2.AdditionalTargets)
}

// CreateProxyTarget ...
func CreateProxyTarget(targetStr string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetStr)
	if err != nil {
		return nil, err
	}
	target := httputil.NewSingleHostReverseProxy(url)
	return target, nil
}
