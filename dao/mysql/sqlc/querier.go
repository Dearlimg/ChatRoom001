// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"
)

type Querier interface {
	CreateAccount(ctx context.Context, arg *CreateAccountParams) error
}

var _ Querier = (*Queries)(nil)
