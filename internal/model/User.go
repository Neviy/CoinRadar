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
	Birthday    string
	Premium     bool
	Admin       bool
	CreateAt    time.Time
}

// NewUser конструктор для инициализации. 
func NewUser(telegramID int64,birthday string)(*User,error){
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
		CreateAt: time.Now(),
	},nil
}
// NewUserFromDB  восстановление пользователя из БД.
func NewUserFromDB(id,telegramID int64,email *string,premium,admin bool,birthday string) *User{
	return &User{
		ID: id,
		TelegramID: telegramID,
		Email: email,
		Birthday:birthday ,
		Premium: premium,
		Admin: admin,
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
func AgeRating(birthday string)error{
	layout := "2006-01-02"
	parsedTime,err:=time.Parse(layout,birthday)
	if err != nil{
		return err
	}
	dateNow:=time.Now()
	duration:=dateNow.Sub(parsedTime)
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
