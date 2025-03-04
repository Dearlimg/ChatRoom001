// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createAccountStmt, err = db.PrepareContext(ctx, createAccount); err != nil {
		return nil, fmt.Errorf("error preparing query CreateAccount: %w", err)
	}
	if q.createUserStmt, err = db.PrepareContext(ctx, createUser); err != nil {
		return nil, fmt.Errorf("error preparing query CreateUser: %w", err)
	}
	if q.deleteUserStmt, err = db.PrepareContext(ctx, deleteUser); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUser: %w", err)
	}
	if q.existEmailStmt, err = db.PrepareContext(ctx, existEmail); err != nil {
		return nil, fmt.Errorf("error preparing query ExistEmail: %w", err)
	}
	if q.existsUserByIDStmt, err = db.PrepareContext(ctx, existsUserByID); err != nil {
		return nil, fmt.Errorf("error preparing query ExistsUserByID: %w", err)
	}
	if q.getAcountIDsByUserIDStmt, err = db.PrepareContext(ctx, getAcountIDsByUserID); err != nil {
		return nil, fmt.Errorf("error preparing query GetAcountIDsByUserID: %w", err)
	}
	if q.getAllEmailStmt, err = db.PrepareContext(ctx, getAllEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllEmail: %w", err)
	}
	if q.getUserByEmailStmt, err = db.PrepareContext(ctx, getUserByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByEmail: %w", err)
	}
	if q.getUserByIDStmt, err = db.PrepareContext(ctx, getUserByID); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserByID: %w", err)
	}
	if q.updateUserStmt, err = db.PrepareContext(ctx, updateUser); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateUser: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createAccountStmt != nil {
		if cerr := q.createAccountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createAccountStmt: %w", cerr)
		}
	}
	if q.createUserStmt != nil {
		if cerr := q.createUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createUserStmt: %w", cerr)
		}
	}
	if q.deleteUserStmt != nil {
		if cerr := q.deleteUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserStmt: %w", cerr)
		}
	}
	if q.existEmailStmt != nil {
		if cerr := q.existEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing existEmailStmt: %w", cerr)
		}
	}
	if q.existsUserByIDStmt != nil {
		if cerr := q.existsUserByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing existsUserByIDStmt: %w", cerr)
		}
	}
	if q.getAcountIDsByUserIDStmt != nil {
		if cerr := q.getAcountIDsByUserIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAcountIDsByUserIDStmt: %w", cerr)
		}
	}
	if q.getAllEmailStmt != nil {
		if cerr := q.getAllEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllEmailStmt: %w", cerr)
		}
	}
	if q.getUserByEmailStmt != nil {
		if cerr := q.getUserByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByEmailStmt: %w", cerr)
		}
	}
	if q.getUserByIDStmt != nil {
		if cerr := q.getUserByIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserByIDStmt: %w", cerr)
		}
	}
	if q.updateUserStmt != nil {
		if cerr := q.updateUserStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateUserStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                       DBTX
	tx                       *sql.Tx
	createAccountStmt        *sql.Stmt
	createUserStmt           *sql.Stmt
	deleteUserStmt           *sql.Stmt
	existEmailStmt           *sql.Stmt
	existsUserByIDStmt       *sql.Stmt
	getAcountIDsByUserIDStmt *sql.Stmt
	getAllEmailStmt          *sql.Stmt
	getUserByEmailStmt       *sql.Stmt
	getUserByIDStmt          *sql.Stmt
	updateUserStmt           *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                       tx,
		tx:                       tx,
		createAccountStmt:        q.createAccountStmt,
		createUserStmt:           q.createUserStmt,
		deleteUserStmt:           q.deleteUserStmt,
		existEmailStmt:           q.existEmailStmt,
		existsUserByIDStmt:       q.existsUserByIDStmt,
		getAcountIDsByUserIDStmt: q.getAcountIDsByUserIDStmt,
		getAllEmailStmt:          q.getAllEmailStmt,
		getUserByEmailStmt:       q.getUserByEmailStmt,
		getUserByIDStmt:          q.getUserByIDStmt,
		updateUserStmt:           q.updateUserStmt,
	}
}
