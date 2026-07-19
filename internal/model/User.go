// package model хранит в себе  сущности

// Здесь хранится сущность User
package model

import (
	"errors"
	"net/mail"
	"time"
)

// User представляет пользователя Telegram-бота.
type User struct{
	ID          int64     
	TelegramID  int64     
	Email       *string 
	Birthday    time.Time
	Premium     bool
	Admin       bool
	CreatedAt    time.Time
}

// NewUser конструктор для инициализации. 
func NewUser(telegramID int64,birthday time.Time)(*User,error){
	if telegramID <= 0 {
		return nil,errors.New("telegram id must be positive")
	}
	if err:=AgeRating(birthday);err != nil{
		return nil,err
	}
	return &User{
		TelegramID: telegramID,
		Email: nil,
		Birthday: birthday,
		Premium: false,
		Admin: false,
		CreatedAt: time.Now(),
	},nil
}
// NewUserFromDB восстанавливает пользователя из БД.
func NewUserFromDB(id int64,
telegramID int64,
	email *string,
	birthday time.Time,
	premium bool,
	admin bool,
	createdAt time.Time,
) *User {
	return &User{
		ID:         id,
		TelegramID: telegramID,
		Email:      email,
		Birthday:   birthday,
		Premium:    premium,
		Admin:      admin,
		CreatedAt:  createdAt,
	}
}

// SetEmail устанавливает адрес электронной почты.
func (u *User)SetEmail(email string)error{
	if _,err:=mail.ParseAddress(email);err!=nil{
		return errors.New("invalid email")
	}
	u.Email=&email
	return nil
}

// AgeRating возрастной цензор
func AgeRating(birthday time.Time)error{
	dateNow:=time.Now()
	duration:=dateNow.Sub(birthday)
	year:=int(duration.Hours()/8766)
	if year >=18{
		return nil
	}
	return errors.New("invalid age")
}

// IsPremium проверяет, является ли пользователь премиум-пользователем.
func (u *User)IsPremium()bool{
	return u.Premium
}

// IsAdmin проверяет, является ли пользователь администратором.
func (u *User)IsAdmin()bool{
	return u.Admin
}

// IsEmailSet проверяет, установлен ли адрес электронной почты.
func (u *User)IsEmailSet()bool{
	return u.Email != nil
}

//MakePremium делает пользователя премиум-пользователем.
func (u *User)MakePremium(){
	u.Premium=true
}

//RemovePremium удаляет премиум-пользователя.
func (u *User)RemovePremium(){
	u.Premium=false
}

//PromoteToAdmin делает пользователя администратором.
func (u *User)PromoteToAdmin(){
	u.Admin=true
}

//PromoteToUser делает пользователя обычным пользователем.
func (u *User)PromoteToUser(){
	u.Admin=false
}
