package ercodes

import "x-bank-ms-bank/cerrors"

const (
	_ cerrors.Code = -iota

	RandomGeneration
	BcryptHashing
	HS512Authorization
	RS256Authorization
)
