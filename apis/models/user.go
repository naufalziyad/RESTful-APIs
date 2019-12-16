package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	FullName  string    `gorm:"size:255;not null; unique" json:"fullname"`
	UserName  string    `gorm:"size:255;not null; unique" json:"username"`
	Email     string    `gorm:"size:100;not null; unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatet_at`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerfyPasswordUser(hasedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hasedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hasedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hasedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.FullName = html.EscapeString(strings.TrimSpace(u.FullName))
	u.UserName = html.EscapeString(strings.TrimSpace(u.UserName))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.FullName == "" {
			return errors.New("Please your Full Name")
		}
		if u.UserName == "" {
			return errors.New("Please your User Name")
		}
		if u.Email == "" {
			return errors.New("Please your Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil

	case "login":
		if u.Password == "" {
			return errors.New("Please input your password")
		}
		if u.Email == "" {
			return errors.New("Please input your email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		//wait tommorow
		return nil
	}
}
