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
	mux.HandleFunc("POST /v1/accounts/block", userMiddlewareGroup.Apply(t.handlerBlockAccount))
	mux.HandleFunc("GET /v1/accounts/history", userMiddlewareGroup.Apply(t.handlerAccountHistory))

	mux.HandleFunc("POST /v1/transactions", userMiddlewareGroup.Apply(t.handlerAccountTransaction))
	mux.HandleFunc("POST /v1/atm/supplement", defaultMiddlewareGroup.Apply(t.handlerATMSupplement))
	mux.HandleFunc("POST /v1/atm/withdrawal", defaultMiddlewareGroup.Apply(t.handlerATMWithdrawal))
	mux.HandleFunc("POST /v1/atm/user/supplement", defaultMiddlewareGroup.Apply(t.handlerATMUserSupplement))

	return mux
}
