package domain

import "errors"

var (

	ErrInvalidInput = errors.New("недействительный_ввод")

	ErrValidationFailed = errors.New("валидация_не_удалась")
)


