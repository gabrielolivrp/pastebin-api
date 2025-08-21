package database

import (
	"errors"

	"gorm.io/gorm"
)

func ErrNotFound(err error) bool {
	if err == nil {
		return true
	}
	return errors.Is(err, gorm.ErrRecordNotFound)
}
