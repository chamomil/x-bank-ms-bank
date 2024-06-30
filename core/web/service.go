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

func (s *Service) OpenAccount(ctx context.Context, userId int64) error {
	return s.accountStorage.OpenUserAccount(ctx, userId)
}

func (s *Service) BlockAccount(ctx context.Context, accountId int64) error {
	return s.accountStorage.BlockUserAccount(ctx, accountId)
}
