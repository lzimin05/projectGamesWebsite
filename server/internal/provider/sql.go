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

func (p *Provider) GetPasswordByEmail(email string) (string, error) {
	var msg string
	err := p.conn.QueryRow("SELECT password FROM users WHERE users.email = ($1)", email).Scan(&msg)
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
