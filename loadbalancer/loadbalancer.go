package loadbalancer

import (
	"cclb/contracts"
	"cclb/log"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type LB struct {
	logger      log.Logger
	mux         *http.ServeMux
	config      Config
	middlewares []func(http.Handler) http.Handler
}

func (rv *LB) RegisterMiddleware(mw func(handler http.Handler) http.Handler) {
	rv.middlewares = append(rv.middlewares, mw)
}

func (rv *LB) RegisterRoute(route url.URL, tg contracts.TargetGroup) {
	rv.logger.Info("Adding routing rule", map[string]any{
		"url":         route.String(),
		"targetGroup": tg.Name(),
	})

	rv.mux.HandleFunc(route.String()+"/", func(w http.ResponseWriter, r *http.Request) {
		rv.logger.Info("Routing via target group", map[string]any{
			"target-group": tg.Name(),
		})

		target := tg.GetNextAvailableTarget()
		target.Serve(w, r)
	})

}

func (rv *LB) Listen() {
	rv.logger.Info("Starting load balancer", nil)

	var handler http.Handler = rv.mux
	for _, mw := range rv.middlewares {
		handler = mw(handler)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", rv.config.Port), handler); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
}

func NewLB(logger log.Logger, config Config) *LB {
	lb := &LB{
		logger: logger,
		config: config,
		mux:    http.NewServeMux(),
	}

	lb.RegisterMiddleware(loggingMiddleware(logger))
	lb.RegisterMiddleware(timeoutMiddleware(logger, config.Timeout))

	return lb
}
