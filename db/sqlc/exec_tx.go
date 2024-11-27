package db

import (
	"context"
	"fmt"
)

// execTX execute a funtion within a database action
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	//begin a database action using the connection pool with the SQLStore
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}

	//creates a new instance of the tx
	q := New(tx)

	//calls the provided function (fn) and passes the value of q
	err = fn(q)
	if err != nil {
		//if is there any errors undo database changes or actions
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	//if is there no errors make database changes or actions
	return tx.Commit(ctx)
}
