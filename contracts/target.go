package contracts

import "net/http"

type Target interface {
	Serve(response http.ResponseWriter, r *http.Request)
}
