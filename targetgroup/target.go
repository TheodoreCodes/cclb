package targetgroup

import (
	"cclb/log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Target struct {
	logger            log.Logger
	URL               url.URL
	ActiveConnections int
	IsHealthy         bool
	mutex             sync.Mutex
}

func (rv *Target) Serve(w http.ResponseWriter, r *http.Request) {
	rv.mutex.Lock()
	rv.ActiveConnections++
	rv.mutex.Unlock()

	proxy := httputil.NewSingleHostReverseProxy(&rv.URL)
	proxy.ErrorHandler = rv.proxyErrHandler()
	proxy.Director = rv.proxyRequestModifier(proxy.Director)

	proxy.ServeHTTP(w, r)

	rv.mutex.Lock()
	rv.ActiveConnections--
	rv.mutex.Unlock()
}

func (rv *Target) proxyRequestModifier(director func(r *http.Request)) func(request *http.Request) {
	return func(r *http.Request) {
		director(r)
		r.URL = &rv.URL
		r.Host = rv.URL.Host

		if _, ok := r.Header["User-Agent"]; !ok {
			r.Header.Set("User-Agent", "")
		}
	}
}

func (rv *Target) proxyErrHandler() func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		rv.logger.Err("failed to reverse proxy", err, map[string]any{
			"target":    rv.URL.String(),
			"isHealthy": rv.IsHealthy,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}
}
