package config

import "net/url"

type TargetGroup struct {
	N       string      `json:"name"`
	InPort  int         `json:"ingressPort"`
	HpInt   int         `json:"healthCheckInterval"`
	HpUri   configUrl   `json:"healthCheckEndpoint"`
	TgtUrls []configUrl `json:"targets"`
}

func (rv *TargetGroup) Name() string {
	return rv.N
}

func (rv *TargetGroup) IngressPort() int {
	return rv.InPort
}

func (rv *TargetGroup) HealthcheckEndpoint() url.URL {
	return url.URL(rv.HpUri)
}

func (rv *TargetGroup) HealthcheckInterval() int {
	return rv.HpInt
}

func (rv *TargetGroup) Targets() []url.URL {
	var res []url.URL
	for idx := range rv.TgtUrls {
		res = append(res, url.URL(rv.TgtUrls[idx]))
	}

	return res
}
