package web

import "time"

type (
	UserAccountsData struct {
		Id           int64
		BalanceCents int64
		Status       string
	}

	AccountTransactionsData struct {
		SenderId    int64
		ReceiverId  int64
		Status      string
		CreatedAt   time.Time
		AmountCents int64
		Description string
	}
)
