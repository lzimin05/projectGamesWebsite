package usecase

import (
	newuser "github.com/lzimin05/IDZ/internal/user"
)

type Provider interface {
	SelectUserByEmail(email string) (string, error)
	SelectUserById(id int) (string, error)
	SelectEmailById(id int) (string, error)
	InsertUser(newUser newuser.User) error
	GetPasswordByEmail(email string) (string, error)
	GetIdByemail(email string) (int, error)
	UpdateSesion(id int) error
	SelectSesion() (int, error)
	UpdateUserById(name string, id int) error
}
