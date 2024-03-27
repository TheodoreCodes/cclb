package main

import (
	"encoding/json"
	"net/url"
)

type Config struct {
	Listeners []ListenerConfig `json:"listeners"`
}

type ListenerConfig struct {
	RoutingUrl        ConfigUrl         `json:"url"`
	TargetGroupConfig TargetGroupConfig `json:"targetGroup"`
}

type TargetGroupConfig struct {
	Name                string      `json:"name"`
	IngressPort         int         `json:"ingressPort"`
	HealthCheckInterval int         `json:"healthCheckInterval"`
	HealthCheckEndpoint ConfigUrl   `json:"healthCheckEndpoint"`
	Targets             []ConfigUrl `json:"targets"`
}

type ConfigUrl url.URL

func (rv *ConfigUrl) UnmarshalJSON(data []byte) error {
	var rawUrl string

	err := json.Unmarshal(data, &rawUrl)
	if err != nil {
		return err
	}

	u, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	*rv = ConfigUrl(*u)

	return nil
}

func (rv *ConfigUrl) URL() *url.URL {
	u := url.URL(*rv)
	return &u
}
