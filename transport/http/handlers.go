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
	var accountData AccountData
	if err := json.NewDecoder(r.Body).Decode(&accountData); err != nil {
		t.errorHandler.setBadRequestError(w, err)
		return
	}

	if err := t.service.BlockAccount(r.Context(), accountData.AccountId); err != nil {
		t.errorHandler.setError(w, err)
	}
	w.WriteHeader(http.StatusOK)
}

func (t *Transport) handlerAccountHistory(w http.ResponseWriter, r *http.Request) {
	var accountData AccountData
	if err := json.NewDecoder(r.Body).Decode(&accountData); err != nil {
		t.errorHandler.setBadRequestError(w, err)
		return
	}

	data, err := t.service.AccountHistory(r.Context(), accountData.AccountId)
	if err != nil {
		t.errorHandler.setError(w, err)
		return
	}

	var response AccountsHistoryResponse
	if data != nil {
		for _, entry := range data {
			userAccountsItem := AccountsHistoryResponseItem{
				SenderId:    entry.SenderId,
				ReceiverId:  entry.ReceiverId,
				Status:      entry.Status,
				CreatedAt:   entry.CreatedAt.Format("2006.01.02 15:04:05"),
				AmountCents: entry.AmountCents,
				Description: entry.Description,
			}
			response.Items = append(response.Items, userAccountsItem)
		}
	} else {
		response.Items = nil
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		t.errorHandler.setError(w, err)
		return
	}
}

func (t *Transport) handlerAccountTransaction(w http.ResponseWriter, r *http.Request) {
	var transactionData TransactionData
	if err := json.NewDecoder(r.Body).Decode(&transactionData); err != nil {
		t.errorHandler.setBadRequestError(w, err)
		return
	}

	if err := t.service.Transaction(r.Context(), transactionData.SenderId, transactionData.ReceiverId, transactionData.AmountCents, transactionData.Description); err != nil {
		t.errorHandler.setError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (t *Transport) handlerATMSupplement(w http.ResponseWriter, r *http.Request) {
	var atmSupplementData ATMOperationData
	if err := json.NewDecoder(r.Body).Decode(&atmSupplementData); err != nil {
		t.errorHandler.setBadRequestError(w, err)
		return
	}

	if err := t.service.ATMSupplement(r.Context(), atmSupplementData.Login, atmSupplementData.Password, atmSupplementData.AmountCents); err != nil {
		t.errorHandler.setError(w, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (t *Transport) handlerATMWithdrawal(w http.ResponseWriter, r *http.Request) {
	var atmWithdrawalData ATMOperationData
	if err := json.NewDecoder(r.Body).Decode(&atmWithdrawalData); err != nil {
		t.errorHandler.setBadRequestError(w, err)
		return
	}

	if err := t.service.ATMWithdrawal(r.Context(), atmWithdrawalData.Login, atmWithdrawalData.Password, atmWithdrawalData.AmountCents); err != nil {
		t.errorHandler.setError(w, err)
	}

	w.WriteHeader(http.StatusOK)
}
