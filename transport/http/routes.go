package http

import "net/http"

func (t *Transport) routes() http.Handler {
	corsHandler := t.corsHandler("*", "*", "*", "")
	corsMiddleware := t.corsMiddleware(corsHandler)

	defaultMiddlewareGroup := middlewareGroup{
		t.panicMiddleware,
		corsMiddleware,
	}

	userMiddlewareGroup := middlewareGroup{
		t.panicMiddleware,
		corsMiddleware,
		t.authMiddleware(false),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", defaultMiddlewareGroup.Apply(t.handlerNotFound))
	mux.HandleFunc("GET /v1/me/accounts", userMiddlewareGroup.Apply(t.handlerUserAccounts))

	mux.HandleFunc("POST /v1/accounts/open", userMiddlewareGroup.Apply(t.handlerOpenAccount))

	return mux
}
