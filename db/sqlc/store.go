package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

/*
 Database Transaction : ACID property
 (A) Atomicity :
 - Either all operations complete successfully or the transaction fails and the db
  is unchanged
(C) Consistency
 - The database state must vaild after the transaction.
   All constrain must be satisfied
 - In another word, all data being written in the db, must not violate any constrain
 (I) Isolation
 - Concurrent transaction must not affect each other
 (D) Durability
 - Data written by a succesful transaction must be recoreded in the persitent storage
*/

// execTx exercutes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q) // call the input function with that qurey
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input paramater of the tranfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Account       int64 `json:"account"`
}

/*
  Create a tranfer record with amount = 10
  Create an account entry for account 1 with amount = - 10
  Create an account entry for account 2 with amount = + 10
  Subtract 10 from the balance of account 1
  Add 10 to the balance of account 2
*/

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	// Entry record money is moving out
	FromEntry Entry `json:"from_entry"`
	// Entry record money is moving in
	ToEntry Entry `json:"to_entry"`
}

// TranferTx perform a money tranfer from one account to another
// It create a tranfer record, add account entries, and update a account balace
// with in a single databse transaction
func (store *Store) TranferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
