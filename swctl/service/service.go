package service

import (
	"github.com/urfave/cli"
)

type service struct {
	flag *cli.Context
	list bool
}

func NewService(flag *cli.Context) *service {
	return &service{
		flag: flag,
		list: flag.Bool("list"),
	}
}

func (s *service) Exec() (err error) {
	if s.list {
		err = s.showList()
	}
	return
}
