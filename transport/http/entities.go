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

	AccountDataToBlock struct {
		AccountId int64
	}
)
