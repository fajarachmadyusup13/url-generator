package usecase

import "errors"

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrURLAlreadyExists = errors.New("url already exists")
)
