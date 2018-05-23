package models

import (
	"errors"
)

var (
	EmptyMessageErr error = errors.New("Can not send empty message")
)
