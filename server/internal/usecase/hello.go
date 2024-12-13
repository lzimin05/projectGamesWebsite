package usecase

import (
	"log"

	newuser "github.com/lzimin05/IDZ/internal/user"
)

func (u *Usecase) PrintUserById(id int) (string, error) {
	msg, err := u.p.SelectUserById(id)
	if err != nil {
		return "", err
	}

	if msg == "" {
		return u.defaultMsg, nil
	}

	return msg, nil
}

func (u *Usecase) PrintUserByEmail(email string) (string, error) {
	msg, err := u.p.SelectUserByEmail(email)
	if err != nil {
		return "", err
	}

	if msg == "" {
		return u.defaultMsg, nil
	}

	return msg, nil
}

func (u *Usecase) InsertNewUser(newUser newuser.User) error {
	flag, err := u.NonUserExistence(newUser)
	if err != nil {
		return err
	}
	if flag {
		//new user registration
		err = u.p.InsertUser(newUser)
		if err != nil {
			return err
		}
		return nil
	}
	log.Println("Пользователь уже зарегистрирован!")
	return nil
}

func (u *Usecase) NonUserExistence(newUser newuser.User) (bool, error) {
	msg, err := u.p.SelectUserByEmail(newUser.Email)
	if err != nil {
		return false, err
	}
	if msg == "" {
		return true, nil
	}
	return false, nil
}

/*
func (u *Usecase) SetHelloMessage(msg string) error {
	isExist, err := u.p.CheckHelloExitByMsg(msg)
	if err != nil {
		return err
	}

	if isExist {
		return nil
	}

	err = u.p.InsertHello(msg)
	if err != nil {
		return err
	}

	return nil
}
*/
