package models

import (
    "github.com/go-playground/validator/v10"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

var validate = validator.New()

func ValidateUser(user *User) error {
    return validate.Struct(user)
}
