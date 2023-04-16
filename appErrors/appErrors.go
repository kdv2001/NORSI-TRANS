package appErrors

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	Msg         string `example:"description"`
	Code        int    `example:"401"`
	Base        error  `json:"-"`
	Description string `json:"-"`
}

var (
	ErrBaseApp    = AppError{"internal server error", http.StatusInternalServerError, nil, ""}
	ErrBadRequest = AppError{"bad request", http.StatusBadRequest, nil, ""}
	ErrForbidden  = AppError{"forbidden", http.StatusForbidden, nil, ""}
	ErrNotFound   = AppError{"not found", http.StatusNotFound, nil, ""}
)

func AppErrorFromError(inputError error) AppError {
	appErr, ok := inputError.(AppError)
	if !ok {
		return ErrBaseApp.Wrap(inputError, "")
	}

	return appErr
}

func (err AppError) IsInternalError() bool {
	return err.Code/100 == 5 //nolint:gomnd
}

func (err AppError) Wrap(baseErr error, desc string) AppError {
	err.Base = baseErr
	err.Description = desc

	return err
}

func (err AppError) Is(target error) bool {
	targetAppErr, ok := target.(AppError)
	if !ok {
		return target == err.Base
	}

	return targetAppErr.Code == err.Code && targetAppErr.Msg == err.Msg
}

func (err AppError) Error() string {
	return err.Msg
}

func (err AppError) String() string {
	errBuffer, er := json.Marshal(err)
	if er != nil {
		panic(er)
	}

	return string(errBuffer)
}
