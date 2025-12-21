package factory

import "errors"

var (
	ErrNotEnoughMoney = errors.New("more coal needed")
	ErrMinerNotExist  = errors.New("miner not exist")
)
