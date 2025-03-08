// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"
)

type Querier interface {
	CountAccountByUserID(ctx context.Context, userID int64) (int64, error)
	CreateAccount(ctx context.Context, arg *CreateAccountParams) error
	CreateUser(ctx context.Context, arg *CreateUserParams) error
	DeleteAccount(ctx context.Context, id int64) error
	DeleteAccountByUserID(ctx context.Context, userID int64) error
	DeleteUser(ctx context.Context, id int64) error
	ExistAccountByID(ctx context.Context, id int64) (bool, error)
	ExistEmail(ctx context.Context, email string) (bool, error)
	ExistsUserByID(ctx context.Context, id int64) (bool, error)
	GetAccountByID(ctx context.Context, arg *GetAccountByIDParams) ([]*GetAccountByIDRow, error)
	GetAccountByUserID(ctx context.Context, userID int64) ([]*GetAccountByUserIDRow, error)
	GetAccountsByName(ctx context.Context, arg *GetAccountsByNameParams) ([]*GetAccountsByNameRow, error)
	GetAcountIDsByUserID(ctx context.Context, userID int64) ([]int64, error)
	GetAllEmail(ctx context.Context) ([]string, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	UpdateAccount(ctx context.Context, arg *UpdateAccountParams) error
	UpdateAccountAvatar(ctx context.Context, arg *UpdateAccountAvatarParams) error
	UpdateUser(ctx context.Context, arg *UpdateUserParams) error
}

var _ Querier = (*Queries)(nil)
