package model

import "errors"

var (
	ErrObjectNotFound = errors.New("not found")
	ErrUpdateFailed   = errors.New("update failed")
	ErrDeleteFailed   = errors.New("delete failed")

	ErrSendReq = errors.New("error with sending request")
)
