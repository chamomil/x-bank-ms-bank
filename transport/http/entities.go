package http

type (
	UserAccountsResponseItem struct {
		Id           int64  `json:"id"`
		BalanceCents int64  `json:"balanceCents"`
		Status       string `json:"status"`
	}

	UserAccountsResponse struct {
		Items []UserAccountsResponseItem `json:"items"`
	}

	AccountData struct {
		AccountId int64 `json:"accountId"`
	}

	AccountsHistoryResponseItem struct {
		SenderId    int64  `json:"senderId"`
		ReceiverId  int64  `json:"receiverId"`
		Status      string `json:"status"`
		CreatedAt   string `json:"createdAt"`
		AmountCents int64  `json:"amountCents"`
		Description string `json:"description"`
	}

	AccountsHistoryResponse struct {
		Items []AccountsHistoryResponseItem `json:"items"`
	}

	TransactionData struct {
		SenderId    int64  `json:"senderId"`
		ReceiverId  int64  `json:"receiverId"`
		AmountCents int64  `json:"amountCents"`
		Description string `json:"description"`
	}

	ATMSupplementData struct {
		Login       string `json:"login"`
		Password    string `json:"password"`
		AmountCents int64  `json:"amountCents"`
	}
)
