package web

import "context"

type (
	AccountStorage interface {
		GetUserAccounts(ctx context.Context, userId int64) ([]UserAccountsData, error)
		OpenUserAccount(ctx context.Context, userId int64) error
		BlockUserAccount(ctx context.Context, accountId int64) error
		GetAccountHistory(ctx context.Context, accountId int64) ([]AccountTransactionsData, error)
	}
)
