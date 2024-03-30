package e

import "fmt"

// Всегда возвращает ненулевую ошибку
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// Возваращает нулевую ошибку если err == nil и вызывает Wrap() если нет
func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}

	return Wrap(msg, err)
}
