package mErrors

import "errors"

var (
	ErrOrderAlreadyUploaded = errors.New("order already uploaded by another user")
	ErrOrderAlreadyCreated  = errors.New("order already created by the same user")
	ErrOrderNotFound        = errors.New("order not found")
	ErrNotEnoughMoney       = errors.New("not enough money")
)
