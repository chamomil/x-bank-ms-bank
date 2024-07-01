package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"x-bank-ms-bank/cerrors"
	"x-bank-ms-bank/core/web"
	"x-bank-ms-bank/ercodes"
)

type (
	Service struct {
		db *sql.DB
	}
)

func NewService(login, password, host string, port int, database string, maxCons int) (Service, error) {
	db, err := sql.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%d/%s", login, password, host, port, database))
	if err != nil {
		return Service{}, err
	}

	db.SetMaxOpenConns(maxCons)

	if err = db.Ping(); err != nil {
		return Service{}, err
	}

	return Service{
		db: db,
	}, err
}

func (s *Service) GetUserAccounts(ctx context.Context, userId int64) ([]web.UserAccountData, error) {
	const query = `SELECT accounts."id", "balanceCents", "status" FROM accounts LEFT JOIN "accountOwners" ON "ownerId" = "accountOwners".id WHERE "userId" = $1`

	rows, err := s.db.QueryContext(ctx, query, userId)

	if err != nil {
		return nil, s.wrapQueryError(err)
	}

	var userAccountsData []web.UserAccountData
	for rows.Next() {
		var data web.UserAccountData
		if err = rows.Scan(&data.Id, &data.BalanceCents, &data.Status); err != nil {
			return nil, s.wrapScanError(err)
		}
		userAccountsData = append(userAccountsData, data)
	}

	return userAccountsData, nil
}

func (s *Service) OpenUserAccount(ctx context.Context, userId int64) error {
	const query = `SELECT "id" FROM "accountOwners" WHERE "userId" = $1`

	row := s.db.QueryRowContext(ctx, query, userId)
	if err := row.Err(); err != nil {
		return s.wrapQueryError(err)
	}

	var accountOwnerId int64
	if err := row.Scan(&accountOwnerId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			accountOwnerId, err = s.createAccountOwner(ctx, userId)
			if err != nil {
				return err
			}
		} else {
			return s.wrapScanError(err)
		}
	}

	const openAccountQuery = `INSERT INTO accounts ("ownerId") VALUES ($1)`
	_, err := s.db.ExecContext(ctx, openAccountQuery, accountOwnerId)
	if err != nil {
		return s.wrapQueryError(err)
	}
	return nil
}

func (s *Service) createAccountOwner(ctx context.Context, userId int64) (int64, error) {
	const query = `INSERT INTO "accountOwners" ("userId") VALUES ($1) RETURNING id`

	row := s.db.QueryRowContext(ctx, query, userId)
	if err := row.Err(); err != nil {
		return 0, s.wrapQueryError(err)
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, s.wrapScanError(err)
	}
	return id, nil
}

func (s *Service) BlockUserAccount(ctx context.Context, accountId int64) error {
	const query = `UPDATE accounts SET status = 'BLOCKED' WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, accountId)
	if err != nil {
		return s.wrapQueryError(err)
	}
	return nil
}

func (s *Service) GetAccountHistory(ctx context.Context, accountId int64) ([]web.AccountTransactionsData, error) {
	const query = `SELECT "senderId", "receiverId", "status", "createdAt", "amountCents", "description" FROM transactions WHERE "senderId" = $1 OR "receiverId" = $1 ORDER BY "createdAt" DESC`

	rows, err := s.db.QueryContext(ctx, query, accountId)
	if err != nil {
		return nil, s.wrapQueryError(err)
	}

	var accountTransactionsData []web.AccountTransactionsData
	for rows.Next() {
		var data web.AccountTransactionsData
		if err = rows.Scan(&data.SenderId, &data.ReceiverId, &data.Status, &data.CreatedAt, &data.AmountCents, &data.Description); err != nil {
			return nil, s.wrapScanError(err)
		}
		accountTransactionsData = append(accountTransactionsData, data)
	}

	return accountTransactionsData, nil
}

func (s *Service) CreateTransaction(ctx context.Context, senderId, receiverId, amountCents int64, description string) error {
	const query = `SELECT "balanceCents", "status" FROM accounts WHERE "id" = $1`
	row := s.db.QueryRowContext(ctx, query, senderId)
	if err := row.Err(); err != nil {
		return s.wrapQueryError(err)
	}

	var userAccountData web.UserAccountData
	if err := row.Scan(&userAccountData.BalanceCents, &userAccountData.Status); err != nil {
		return s.wrapScanError(err)
	}
	if userAccountData.Status == "BLOCKED" {
		return cerrors.NewErrorWithUserMessage(ercodes.BlockedAccount, nil, "Счёт отправителя заблокирован")
	}
	if userAccountData.BalanceCents < amountCents {
		return cerrors.NewErrorWithUserMessage(ercodes.NotEnoughMoney, nil, "Недостаточно средств")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return s.wrapQueryError(err)
	}
	defer tx.Rollback()
	const queryTransaction = `INSERT INTO transactions ("senderId", "receiverId", "amountCents", description) VALUES (@senderId, @receiverId, @amountCents, @description)`
	_, err = tx.ExecContext(ctx, queryTransaction, pgx.NamedArgs{
		"senderId":    senderId,
		"receiverId":  receiverId,
		"amountCents": amountCents,
		"description": description,
	})
	if err != nil {
		return s.wrapQueryError(err)
	}

	const querySenderUpdate = `UPDATE accounts SET "balanceCents" = "balanceCents" - @amountCents WHERE id = @senderId`
	const queryReceiverUpdate = `UPDATE accounts SET "balanceCents" = "balanceCents" + @amountCents WHERE id = @receiverId`

	_, err = s.db.ExecContext(ctx, querySenderUpdate, pgx.NamedArgs{
		"amountCents": amountCents,
		"senderId":    senderId,
	})
	if err != nil {
		return s.wrapQueryError(err)
	}

	_, err = s.db.ExecContext(ctx, queryReceiverUpdate, pgx.NamedArgs{
		"amountCents": amountCents,
		"receiverId":  receiverId,
	})
	if err != nil {
		return s.wrapQueryError(err)
	}
	if err = tx.Commit(); err != nil {
		return s.wrapQueryError(err)
	}
	return nil
}
