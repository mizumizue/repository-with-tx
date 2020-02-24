package rdb

import "fmt"

type BeginTxErr struct {
	err error
}

func (e BeginTxErr) Error() string {
	return fmt.Sprintf("begin transaction failed. detail: %s", e.err.Error())
}

func (e BeginTxErr) Unwrap() error {
	return e.err
}
