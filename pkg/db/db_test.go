package db

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func TestConstructorAndConnector(t *testing.T) {
	connector := Db{}
	connector = connector.Constructor("sqlite3", ":memory:")
	if connector.connector == nil {
		t.Errorf("connector should contain *sqlx.db")
	}
	// setconnector
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	connector.SetConnector(sqlx.NewDb(db, "sqlmock"))
	if connector.connector == nil {
		t.Errorf("connector should contain *sqlx.db")
	}

	if connector.Connector() == nil {
		t.Errorf("connector should contain *sqlx.db")
	}
}

func TestFindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	connector := Db{}
	connector.connector = sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	mock.ExpectPrepare("^SELECT (.+) FROM alerts xx").ExpectQuery()
	mock.ExpectPrepare("^SELECT (.+) FROM alerts").ExpectQuery().WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectPrepare("^SELECT (.+) FROM alerts").ExpectQuery().WillReturnRows(rows)
	mock.ExpectCommit()

	type Alert struct {
		Id    int
		Title string
		Body  string
	}

	alerts := []Alert{}

	// test with invalidd query
	_, err = connector.FindAll(Params{
		Query: "SELECT id,title FROM alerts xx",
		Dest:  &alerts,
	})

	if err == nil || len(alerts) != 0 {
		t.Errorf("error was expected if sql contains invalid query")
	}

	// test with valid value
	_, err = connector.FindAll(Params{
		Query: "SELECT id,title FROM alerts",
		Dest:  &alerts,
	})

	if err != nil {
		t.Errorf("error was not expected : %s", err)
	}

	if len(alerts) != 2 {
		t.Errorf("data should have 2 rows data")
	}

	// test with tx
	tx, err := connector.Begin()
	if err != nil {
		t.Errorf("error was not expected on initialization exec : %s", err)
	}

	_, err = connector.FindAll(Params{
		Tx:    tx,
		Query: "SELECT id,subject FROM alerts",
		Dest:  &alerts,
	})

	if err != nil {
		t.Errorf("error was not expected on transact exec : %s", err)
	}

	err = tx.Commit()

	if err != nil {
		t.Errorf("error was not expected on transact commit : %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFindOne(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	connector := Db{}
	connector.connector = sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	mock.ExpectPrepare("^SELECT (.+) FROM alerts limit 1").ExpectQuery().WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectPrepare("^SELECT (.+) FROM alerts limit 1").ExpectQuery().WillReturnRows(rows)
	mock.ExpectCommit()

	type Alert struct {
		Id    int
		Title string
		Body  string
	}

	alert := Alert{}

	// test with valid value
	_, err = connector.FindOne(Params{
		Query: "SELECT id,subject FROM alerts limit 1",
		Dest:  &alert,
	})

	if err != nil {
		t.Errorf("error was not expected : %s", err)
	}

	if alert.Body != "hello" {
		t.Errorf("body should contain hello, instead got : %s", alert.Body)
	}

	// test with tx
	tx, err := connector.Begin()
	if err != nil {
		t.Errorf("error was not expected on initialization exec : %s", err)
	}

	_, err = connector.FindOne(Params{
		Tx:    tx,
		Query: "SELECT id,subject FROM alerts limit 1",
		Dest:  &alert,
	})

	if err != nil {
		t.Errorf("error was not expected on transact exec : %s", err)
	}

	err = tx.Commit()

	if err != nil {
		t.Errorf("error was not expected on transact commit : %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	connector := Db{}
	connector.connector = sqlx.NewDb(db, "sqlmock")

	mock.ExpectPrepare("INSERT INTO alerts").ExpectExec().WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO alerts").ExpectExec().WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// test with valid value
	_, _, err = connector.Exec(Params{
		Query: "INSERT INTO alerts (a_id, b_id) VALUES (?, ?)",
		Args:  []interface{}{2, 3},
	})

	if err != nil {
		t.Errorf("error was not expected : %s", err)
	}

	// test with tx
	tx, err := connector.Begin()
	if err != nil {
		t.Errorf("error was not expected on initialization exec : %s", err)
	}

	_, _, err = connector.Exec(Params{
		Tx:    tx,
		Query: "INSERT INTO alerts (a_id, b_id) VALUES (?, ?)",
		Args:  []interface{}{2, 3},
	})

	if err != nil {
		t.Errorf("error was not expected on transact exec : %s", err)
	}

	err = tx.Commit()

	if err != nil {
		t.Errorf("error was not expected on transact commit : %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertAndReturnID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	connector := Db{}
	connector.connector = sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{"id", "title", "body"}).
		AddRow(1, "post 1", "hello").
		AddRow(2, "post 2", "world")

	mock.ExpectPrepare("INSERT INTO alerts").ExpectQuery().WithArgs(2, 3).WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO alerts").ExpectQuery().WithArgs(2, 3).WillReturnRows(rows)
	mock.ExpectCommit()

	// test with valid value
	_, _, err = connector.InsertAndReturnID(Params{
		Query: "INSERT INTO alerts (a_id, b_id) VALUES (?, ?)",
		Args:  []interface{}{2, 3},
	})

	if err != nil {
		t.Errorf("error was not expected : %s", err)
	}

	// test with tx
	tx, err := connector.Begin()
	if err != nil {
		t.Errorf("error was not expected on initialization exec : %s", err)
	}

	_, _, err = connector.InsertAndReturnID(Params{
		Tx:    tx,
		Query: "INSERT INTO alerts (a_id, b_id) VALUES (?, ?)",
		Args:  []interface{}{2, 3},
	})

	if err != nil {
		t.Errorf("error was not expected on transact exec : %s", err)
	}

	err = tx.Commit()

	if err != nil {
		t.Errorf("error was not expected on transact commit : %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
