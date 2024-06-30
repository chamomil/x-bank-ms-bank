package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"x-bank-ms-bank/auth"
)

func (t *Transport) handlerNotFound(w http.ResponseWriter, _ *http.Request) {
	t.errorHandler.setNotFoundError(w)
}

func (t *Transport) handlerUserAccounts(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(t.claimsCtxKey).(*auth.Claims)
	if !ok {
		t.errorHandler.setError(w, errors.New("отсутствуют claims в контексте"))
		return
	}
	userId := claims.Sub
	data, err := t.service.GetAccounts(r.Context(), userId)
	if err != nil {
		t.errorHandler.setError(w, err)
		return
	}

	var response UserAccountsResponse
	if data != nil {
		for _, entry := range data {
			userAccountsItem := UserAccountsResponseItem{
				Id:           entry.Id,
				BalanceCents: entry.BalanceCents,
				Status:       entry.Status,
			}
			response.Items = append(response.Items, userAccountsItem)
		}
	} else {
		response.Items = nil
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		t.errorHandler.setError(w, err)
		return
	}
}

func (t *Transport) handlerOpenAccount(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(t.claimsCtxKey).(*auth.Claims)
	if !ok {
		t.errorHandler.setError(w, errors.New("отсутствуют claims в контексте"))
		return
	}
	userId := claims.Sub
	if err := t.service.OpenAccount(r.Context(), userId); err != nil {
		t.errorHandler.setError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (t *Transport) handlerBlockAccount(w http.ResponseWriter, r *http.Request) {
	var accountData AccountDataToBlock
	if err := json.NewDecoder(r.Body).Decode(&accountData); err != nil {
		t.errorHandler.setBadRequestError(w, err)
		return
	}

	if err := t.service.BlockAccount(r.Context(), accountData.AccountId); err != nil {
		t.errorHandler.setError(w, err)
	}
	w.WriteHeader(http.StatusOK)
}
