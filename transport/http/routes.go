package http

import "net/http"

func (t *Transport) routes() http.Handler {
	mux := http.NewServeMux()
	return mux
}
