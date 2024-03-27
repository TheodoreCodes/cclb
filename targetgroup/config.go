package targetgroup

import "net/url"

type Config interface {
	Name() string
	IngressPort() int
	HealthcheckEndpoint() url.URL
	HealthcheckInterval() int
	Targets() []url.URL
}
