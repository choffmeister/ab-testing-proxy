package main

import (
	"testing"
)

func TestLoadABTestingProxyConfig(t *testing.T) {
	yaml := `cookieName: "test"
defaultTarget:
  id: "default"
  url: "http://localhost:10000"
additionalTargets:
- id: other
  url: "http://localhost:10001"
`

	c1, _, _, err := LoadABTestingProxyConfig([]byte(yaml))
	if err != nil {
		t.Error(err)
		return
	}
	c2 := ABTestingProxyConfig{
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
	}
	if !c1.Equal(c2) {
		t.Errorf("expected %v, got %v", c2, c1)
		return
	}
}
