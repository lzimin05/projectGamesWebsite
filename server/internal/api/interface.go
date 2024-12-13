package api

import (
	newuser "github.com/lzimin05/IDZ/internal/user"
)

type Usecase interface {
	PrintUserById(id int) (string, error)
	PrintUserByEmail(email string) (string, error)
	InsertNewUser(newUser newuser.User) error
	NonUserExistence(newUser newuser.User) (bool, error)
	GetPasswordByEmail(email string) (string, error)
}
