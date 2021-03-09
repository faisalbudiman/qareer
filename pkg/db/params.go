package db

import (
	"github.com/jmoiron/sqlx"
)

type Params struct {
	Tx    *sqlx.Tx
	Query string
	Args  []interface{}
	Dest  interface{}
}
