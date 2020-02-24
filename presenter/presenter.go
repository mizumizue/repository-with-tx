package presenter

import (
	"io"
	"os"
)

type IPresenter interface {
	io.Writer
}

type StdoutPresenter struct {
	io.Writer
}

func NewStdoutPresenter() IPresenter {
	return &StdoutPresenter{os.Stdout}
}
