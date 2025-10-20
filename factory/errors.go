package factory

import "errors"

var (
	ErrNotEnoughMoney = errors.New("more coal nedeed")
	ErrMinerMotExist  = errors.New("miner not exist")
)
