package web

import "context"

type (
	Service struct {
		accountStorage AccountStorage
	}
)

func NewService(accountStorage AccountStorage) Service {
	return Service{
		accountStorage: accountStorage,
	}
}

func (s *Service) GetAccounts(ctx context.Context, userId int64) ([]UserAccountsData, error) {
	return s.accountStorage.GetUserAccounts(ctx, userId)
}
