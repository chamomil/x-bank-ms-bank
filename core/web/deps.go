package web

import "context"

type (
	AccountStorage interface {
		GetUserAccounts(ctx context.Context, userId int64) ([]UserAccountData, error)
		OpenUserAccount(ctx context.Context, userId int64) error
		BlockUserAccount(ctx context.Context, accountId int64) error
		GetAccountHistory(ctx context.Context, accountId int64) ([]AccountTransactionsData, error)
		UpdateAtmAccount(ctx context.Context, amountCents, accountId int64) error
		GetSenderAccountData(ctx context.Context, senderId int64) (UserAccountData, error)
	}

	TransactionStorage interface {
		CreateTransaction(ctx context.Context, senderId, receiverId, amountCents int64, description string) error
	}

	AtmStorage interface {
		GetPasswordByLogin(ctx context.Context, login string) (AtmData, error)
		UpdateAtmCash(ctx context.Context, amountCents, atmId int64) error
	}

	PasswordHasher interface {
		CompareHashAndPassword(ctx context.Context, password string, hashedPassword []byte) error
		HashPassword(_ context.Context, password []byte, cost int) ([]byte, error)
	}
)
