package targetgroup

import (
	"cclb/contracts"
	"cclb/log"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TargetGroup struct {
	logger              log.Logger
	name                string
	ingressPort         int
	healthCheckInterval int
	healthCheckEndpoint url.URL
	targets             []Target
}

func (rv *TargetGroup) RegisterTarget(addr url.URL) error {
	host := fmt.Sprintf("%s:%d", addr.Host, rv.ingressPort)
	addr.Host = host

	rv.targets = append(rv.targets, Target{
		logger: rv.logger,
		URL:    addr,
	})

	return nil
}

func (rv *TargetGroup) startHealthCheck() {
	rv.logger.Info(fmt.Sprintf("Starting healthcheck routing for target group %s", rv.name), nil)

	for idx := range rv.targets {
		var t = &rv.targets[idx]
		rv.logger.Info("Starting healthcheck", map[string]any{
			"target": t.URL.String(),
		})
		go func(t *Target, healthCheckEndpoint string) {
			for range time.Tick(time.Duration(rv.healthCheckInterval) * time.Second) {
				healthCheckUrl := t.URL
				healthCheckUrl.Path = healthCheckEndpoint
				res, err := http.Get(healthCheckUrl.String())
				t.mutex.Lock()
				if err != nil || res.StatusCode >= 500 {
					t.IsHealthy = false
				} else {
					t.IsHealthy = true
				}
				t.mutex.Unlock()

				rv.logger.Debug("Healthcheck result", map[string]any{
					"target":    t.URL.String(),
					"isHealthy": t.IsHealthy,
				})
			}
		}(t, rv.healthCheckEndpoint.Path)
	}
}

func (rv *TargetGroup) Name() string {
	return rv.name
}

func (rv *TargetGroup) GetNextAvailableTarget() contracts.Target {
	currTarget := &rv.targets[0]
	for i := range rv.targets {
		rv.targets[i].mutex.Lock()
		rv.logger.Debug("isHealthy", map[string]any{
			"target":    rv.targets[i].URL,
			"isHealthy": rv.targets[i].IsHealthy,
		})
		if rv.targets[i].ActiveConnections < currTarget.ActiveConnections && rv.targets[i].IsHealthy {

			currTarget = &rv.targets[i]
		}

		rv.targets[i].mutex.Unlock()
	}

	return currTarget
}

func NewTargetGroup(logger log.Logger, conf Config) (contracts.TargetGroup, error) {
	tg := &TargetGroup{
		logger:              logger,
		name:                conf.Name(),
		ingressPort:         conf.IngressPort(),
		healthCheckEndpoint: conf.HealthcheckEndpoint(),
		healthCheckInterval: conf.HealthcheckInterval(),
	}

	for _, tUrl := range conf.Targets() {
		if err := tg.RegisterTarget(tUrl); err != nil {
			return nil, err
		}
	}

	tg.startHealthCheck()

	return tg, nil
}
