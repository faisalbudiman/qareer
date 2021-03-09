package db

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type iconnector interface {
	Preparex(string) (*sqlx.Stmt, error)
}

type ilogger interface {
	Println(...interface{})
	Errorln(...interface{})
}

type Db struct {
	connector *sqlx.DB
	logger    ilogger
}

type Itx interface {
	Commit() error
	Rollback() error
}

func (db Db) Constructor(driver string, url string) Db {
	db.connector = sqlx.MustConnect(driver, url)
	return db
}

func (db Db) SetConnector(connector *sqlx.DB) Db {
	db.connector = connector
	return db
}

func (db Db) Connector() *sqlx.DB {
	return db.connector
}

func (db Db) Begin() (*sqlx.Tx, error) {
	return db.connector.Beginx()
}

func (db Db) Select(dest interface{}, q string, args ...interface{}) error {
	_, err := db.FindAll(Params{
		Query: q,
		Args:  args,
		Dest:  dest,
	})
	return err
}

func (db Db) FindAll(params Params) (duration time.Duration, err error) {
	now := time.Now()
	var connector iconnector
	if params.Tx != nil {
		connector = params.Tx
	} else {
		connector = db.connector
	}

	stmt, err := connector.Preparex(params.Query)
	if err != nil {
		duration = time.Since(now).Round(time.Microsecond)
		return
	}

	err = stmt.Select(params.Dest, params.Args...)
	if err != nil {
		duration = time.Since(now).Round(time.Microsecond)
		return
	}

	duration = time.Since(now).Round(time.Microsecond)
	return
}

func (db Db) FindOne(params Params) (duration time.Duration, err error) {
	now := time.Now()
	var connector iconnector
	if params.Tx != nil {
		connector = params.Tx
	} else {
		connector = db.connector
	}

	stmt, err := connector.Preparex(params.Query)
	if err != nil {
		duration = time.Since(now).Round(time.Microsecond)
		return
	}

	err = stmt.Get(params.Dest, params.Args...)
	duration = time.Since(now).Round(time.Microsecond)
	return
}

func (db Db) InsertAndReturnID(params Params) (int, time.Duration, error) {
	now := time.Now()
	id := 0
	var connector iconnector
	if params.Tx != nil {
		connector = params.Tx
	} else {
		connector = db.connector
	}

	stmt, err := connector.Preparex(params.Query)
	if err != nil {
		return id, time.Since(now).Round(time.Microsecond), err
	}

	rows, err := stmt.Queryx(params.Args...)
	if err != nil {
		return id, time.Since(now).Round(time.Microsecond), err
	}

	for rows.Next() {
		rows.Scan(&id)
	}

	return id, time.Since(now).Round(time.Microsecond), err
}

func (db Db) Save(q string, args ...interface{}) error {
	_, _, err := db.Exec(Params{
		Query: q,
		Args:  args,
	})
	return err
}

func (db Db) Exec(params Params) (sql.Result, time.Duration, error) {
	var result sql.Result
	now := time.Now()
	var connector iconnector
	if params.Tx != nil {
		connector = params.Tx
	} else {
		connector = db.connector
	}

	stmt, err := connector.Preparex(params.Query)
	if err != nil {
		return result, time.Since(now).Round(time.Microsecond), err
	}

	result, err = stmt.Exec(params.Args...)
	if err != nil {
		return result, time.Since(now).Round(time.Microsecond), err
	}
	return result, time.Since(now).Round(time.Microsecond), err
}
