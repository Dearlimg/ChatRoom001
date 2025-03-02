package operate

import (
	"context"
)

const EmailKey = "EmailKey"

func (r *RDB) AddEmail(ctx context.Context, email ...string) error {
	return nil
}

func (r *RDB) ExistEmail(ctx context.Context, email string) (bool, error) {
	return r.rdb.Get(ctx, EmailKey).Bool()
}
