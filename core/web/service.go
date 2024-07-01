package web

import (
	"context"
	"fmt"
)

type (
	Service struct {
		accountStorage AccountStorage
		passwordHasher PasswordHasher
		atmStorage     AtmStorage
	}
)

func NewService(accountStorage AccountStorage, passwordHasher PasswordHasher, atmStorage AtmStorage) Service {
	return Service{
		accountStorage: accountStorage,
		passwordHasher: passwordHasher,
		atmStorage:     atmStorage,
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

func (s *Service) ATMSupplement(ctx context.Context, login, password string, amountCents int64) error {
	atmData, err := s.atmStorage.GetPasswordByLogin(ctx, login)
	if err != nil {
		return err
	}
	a, _ := s.passwordHasher.HashPassword(ctx, []byte(password), 10)
	fmt.Println(string(a))

	if err = s.passwordHasher.CompareHashAndPassword(ctx, password, atmData.PasswordHash); err != nil {
		return err
	}

	if err = s.atmStorage.UpdateAtmCash(ctx, amountCents, atmData.Id); err != nil {
		return err
	}
	if err = s.accountStorage.UpdateAtmAccount(ctx, amountCents, atmData.AccountId); err != nil {
		return err
	}
	return nil
}
