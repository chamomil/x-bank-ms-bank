package web

import (
	"context"
	"x-bank-ms-bank/cerrors"
	"x-bank-ms-bank/ercodes"
)

type (
	Service struct {
		accountStorage     AccountStorage
		passwordHasher     PasswordHasher
		atmStorage         AtmStorage
		transactionStorage TransactionStorage
	}
)

func NewService(accountStorage AccountStorage, passwordHasher PasswordHasher, atmStorage AtmStorage, transactionStorage TransactionStorage) Service {
	return Service{
		accountStorage:     accountStorage,
		passwordHasher:     passwordHasher,
		atmStorage:         atmStorage,
		transactionStorage: transactionStorage,
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
	senderAccountData, err := s.accountStorage.GetSenderAccountData(ctx, senderId)
	if err != nil {
		return err
	}

	if senderAccountData.Status == "BLOCKED" {
		return cerrors.NewErrorWithUserMessage(ercodes.BlockedAccount, nil, "Счёт отправителя заблокирован")
	}
	if senderAccountData.BalanceCents < amountCents {
		return cerrors.NewErrorWithUserMessage(ercodes.NotEnoughMoney, nil, "Недостаточно средств")
	}

	return s.transactionStorage.CreateTransaction(ctx, senderId, receiverId, amountCents, description)
}

func (s *Service) ATMSupplement(ctx context.Context, login, password string, amountCents int64) error {
	atmData, err := s.atmStorage.GetPasswordByLogin(ctx, login)
	if err != nil {
		return err
	}

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

func (s *Service) ATMWithdrawal(ctx context.Context, login, password string, amountCents int64) error {
	atmData, err := s.atmStorage.GetPasswordByLogin(ctx, login)
	if err != nil {
		return err
	}

	if err = s.passwordHasher.CompareHashAndPassword(ctx, password, atmData.PasswordHash); err != nil {
		return err
	}

	if err = s.atmStorage.UpdateAtmCash(ctx, -amountCents, atmData.Id); err != nil {
		return err
	}
	if err = s.accountStorage.UpdateAtmAccount(ctx, -amountCents, atmData.AccountId); err != nil {
		return err
	}
	return nil
}
