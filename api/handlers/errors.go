package handlers

import "errors"

var (
	ErrCannotParse = errors.New("item must be provided either in JSON body or query parameter")
)
