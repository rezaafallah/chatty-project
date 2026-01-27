package validator

import (
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	once     sync.Once
	validate *validator.Validate
)

func GetInstance() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}

// ValidateStruct
func ValidateStruct(s interface{}) error {
	return GetInstance().Struct(s)
}