package rdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DBConn struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewDBConn() (SqlHandler, error) {
	dbConn, err := sqlx.Connect("mysql", "mysql:password@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	return &DBConn{
		db: dbConn,
	}, nil
}

func (conn *DBConn) InTransaction() bool {
	return conn.tx != nil
}

func (conn *DBConn) Transact(ctx context.Context, txFunc func() error) (err error) {
	tx, err := conn.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		err = BeginTxErr{err: err}
		return
	}
	conn.tx = tx

	defer func() {
		defer func() {
			conn.tx = nil
		}()
		if p := recover(); p != nil {
			err = tx.Rollback()
			if err != nil {
				panic(fmt.Errorf("rollback failed. detail: %v", err))
			}
			panic(p) // re-throw panic after Rollback
		}
		if err != nil && errors.As(err, &BeginTxErr{}) {
			return
		}
		if err != nil {
			txErr := tx.Rollback() // err is nil; if Commit returns error update err
			if txErr != nil {
				err = fmt.Errorf("transaction failed. detail: %w", err)
			}
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = txFunc()
	return err
}

func (conn *DBConn) Close() error {
	return conn.db.Close()
}

func (conn *DBConn) get(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	if conn.InTransaction() {
		return conn.tx.GetContext(ctx, dst, query, args...)
	}
	return conn.tx.GetContext(ctx, dst, query, args...)
}

func (conn *DBConn) query(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	if conn.InTransaction() {
		return conn.tx.QueryxContext(ctx, query, args...)
	}
	return conn.db.QueryxContext(ctx, query, args...)
}

func (conn *DBConn) exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if conn.InTransaction() {
		return conn.tx.ExecContext(ctx, query, args...)
	}
	return conn.db.ExecContext(ctx, query, args...)
}

func (conn *DBConn) namedQuery(ctx context.Context, query string, namedParam map[string]interface{}) (*sqlx.Rows, error) {
	panic("implement me")
}

func (conn *DBConn) namedExec(ctx context.Context, query string, namedParam map[string]interface{}) (sql.Result, error) {
	panic("implement me")
}

func (conn *DBConn) rowsScan(ctx context.Context, rows *sqlx.Rows, typeI interface{}) (interface{}, error) {
	defer func() {
		if rows != nil {
			return
		}
		_ = rows.Close()
	}()

	sliceType := reflect.SliceOf(reflect.PtrTo(reflect.TypeOf(typeI)))
	list := reflect.New(sliceType).Elem()

	for rows.Next() {
		dst := reflect.New(reflect.TypeOf(typeI))
		dstI := dst.Interface()
		if err := rows.StructScan(dstI); err != nil {
			return nil, err
		}
		list = reflect.Append(list, dst)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list.Interface(), nil
}
