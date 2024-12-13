package provider

import (
	"database/sql"
	"errors"

	newuser "github.com/lzimin05/IDZ/internal/user"
)

func (p *Provider) SelectUserById(id int) (string, error) {
	var msg string
	err := p.conn.QueryRow("SELECT name FROM users WHERE users.id = ($1)", id).Scan(&msg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return msg, nil
}

func (p *Provider) SelectUserByEmail(email string) (string, error) {
	var msg string

	err := p.conn.QueryRow("SELECT name FROM users WHERE users.email = ($1)", email).Scan(&msg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return msg, nil
}

func (p *Provider) InsertUser(newUser newuser.User) error {
	_, err := p.conn.Exec("INSERT INTO users(name, email, password) VALUES ($1, $2, $3)", newUser.Name, newUser.Email, newUser.Password)
	if err != nil {
		return err
	}
	return nil
}

/*
func (p *Provider) CheckHelloExitByMsg(msg string) (bool, error) {
	// Получаем одно сообщение из таблицы hello
	err := p.conn.QueryRow("SELECT message FROM hello WHERE message = $1 LIMIT 1", msg).Scan(&msg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (p *Provider) InsertHello(msg string) error {
	_, err := p.conn.Exec("INSERT INTO hello (message) VALUES ($1)", msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) SelectRandomHello() (string, error) {
	var msg string

	// Получаем одно сообщение из таблицы hello, отсортированной в случайном порядке
	err := p.conn.QueryRow("SELECT message FROM hello ORDER BY RANDOM() LIMIT 1").Scan(&msg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	return msg, nil
}

*/
