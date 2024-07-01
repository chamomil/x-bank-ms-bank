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

func (s *Service) GetAccounts(ctx context.Context, userId int64) ([]UserAccountData, error) {
	return s.accountStorage.GetUserAccounts(ctx, userId)
}

func (s *Service) OpenAccount(ctx context.Context, userId int64) error {
	return s.accountStorage.OpenUserAccount(ctx, userId)
}

func (s *Service) BlockAccount(ctx context.Context, accountId int64) error {
	return s.accountStorage.BlockUserAccount(ctx, accountId)
}

func (s *Service) AccountHistory(ctx context.Context, accountId int64) ([]AccountTransactionsData, error) {
	return s.accountStorage.GetAccountHistory(ctx, accountId)
}

func (s *Service) Transaction(ctx context.Context, senderId, receiverId, amountCents int64, description string) error {
	return s.accountStorage.CreateTransaction(ctx, senderId, receiverId, amountCents, description)
}
