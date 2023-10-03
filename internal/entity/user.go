package entity

import "golang.org/x/crypto/bcrypt"

//go:generate easyjson user.go
//easyjson:json
type User struct {
	ID       uint64
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (c *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return err
	}
	c.Password = string(bytes)
	return nil
}

func (c *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password))
}
