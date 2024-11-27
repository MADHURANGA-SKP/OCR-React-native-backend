package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// store defines all funtions to execute db quaries withing db action
type Store interface {
	Querier
}

// store provide all funtions to execute db queries and data trival and transfers
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// create NewStore
func NewStore(connPool *pgxpool.Pool) Store { //pgxpool.Pool manages connections to a  db
	return &SQLStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
