package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"x-bank-ms-bank/core/web"
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

func (s *Service) GetUserAccounts(ctx context.Context, userId int64) ([]web.UserAccountsData, error) {
	const query = `SELECT accounts."id", "balanceCents", "status" FROM accounts LEFT JOIN "accountOwners" ON "ownerId" = "accountOwners".id WHERE "userId" = $1`

	rows, err := s.db.QueryContext(ctx, query, userId)

	if err != nil {
		return nil, s.wrapQueryError(err)
	}

	var userAccountsData []web.UserAccountsData
	for rows.Next() {
		var data web.UserAccountsData
		if err = rows.Scan(&data.Id, &data.BalanceCents, &data.Status); err != nil {
			return nil, s.wrapScanError(err)
		}
		userAccountsData = append(userAccountsData, data)
	}

	return userAccountsData, nil
}
