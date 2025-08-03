package model

import (
	"errors"
	"regexp"
	"unicode"
)

type User struct {
	UUID     string `json:"uuid"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (u *User) ValidateUserName() error {
	for _, r := range u.Name {
		if !unicode.IsLetter(r) {
			return errors.New("Hata isim!")
		}
	}

	return nil
}
func (u *User) ValidateEmail() error {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, u.Email)
	if !matched {
		return errors.New("Hata email!")
	}
	return nil
}
func (u *User) ValidatePassword() error {
	if len(u.Password) <= 3 {
		return errors.New("Hata şifre!")
	}
	return nil
}
func (u *User) ValidateRole() error {
	if u.Role != "user" {
		return errors.New("Hata rol!")
	}
	return nil
}
func (u *User) ValidateAge() error {
	if u.Age < 0 || u.Age >= 100 {
		return errors.New("ya doğmadın ya öldün")
	}
	return nil
}
