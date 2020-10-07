package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

func main() {
	configPath := flag.String("config", "config.yaml", "The config file path.")
	listen := flag.String("listen", "0.0.0.0:8080", "The interface to listen on.")
	flag.Parse()

	config, defaultProxy, additionalProxies, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("unable to load config: %v\n", err)
		return
	}
	proxy := &ABTestingProxy{
		Config:            *config,
		DefaultProxy:      *defaultProxy,
		AdditionalProxies: *additionalProxies,
	}

	go reloadConfigLoop(*configPath, proxy)

	log.Printf("starting proxy server on %s\n", *listen)
	if err := http.ListenAndServe(*listen, proxy); err != nil {
		log.Fatalf("unable to start server on %s: %v\n", *listen, err)
	}
}

func loadConfig(configPath string) (*ABTestingProxyConfig, *httputil.ReverseProxy, *map[string]httputil.ReverseProxy, error) {
	configYaml, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, nil, nil, err
	}
	config, defaultProxy, additionalProxies, err := LoadABTestingProxyConfig(configYaml)
	if err != nil {
		return nil, nil, nil, err
	}
	return config, defaultProxy, additionalProxies, nil
}

func reloadConfigLoop(configPath string, proxy *ABTestingProxy) {
	for {
		time.Sleep(10 * time.Second)
		config, defaultProxy, additionalProxies, err := loadConfig(configPath)
		if err != nil {
			log.Printf("unable to reload config: %v\n", err)
			continue
		}

		if config.Equal(proxy.Config) {
			continue
		}

		proxy.Config = *config
		proxy.DefaultProxy = *defaultProxy
		proxy.AdditionalProxies = *additionalProxies
		log.Printf("reloaded config\n")
	}
}
