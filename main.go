package main

import (
	"cclb/config"
	"cclb/contracts"
	"cclb/loadbalancer"
	"cclb/log"
	"cclb/targetgroup"
	"net/url"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	env, err := LoadEnv()
	if err != nil {
		return err
	}

	logger := log.NewDefaultLogger(env.LogLevel)
	cfg := config.Load(logger, env.ConfigFile)

	lbTimeout, err := time.ParseDuration(env.Timeout)
	if err != nil {
		logger.Err("failed to parse load balancer timeout", err, nil)
		return err
	}

	lb := loadbalancer.NewLB(logger, loadbalancer.Config{
		Port:    env.Port,
		Timeout: lbTimeout,
	})

	var tg contracts.TargetGroup
	for i := range cfg.Listeners {
		tg, _ = targetgroup.NewTargetGroup(
			logger,
			&cfg.Listeners[i].TargetGroupConfig,
		)

		lb.RegisterRoute(url.URL(cfg.Listeners[i].RoutingUrl), tg)
	}

	lb.Listen()
	return nil
}
