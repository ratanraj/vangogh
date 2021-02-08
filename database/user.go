package database

import (
	"fmt"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model

	Email string `json:"email" gorm:"uniqueIndex" form:"email"`
	FirstName string `json:"first_name" form:"first_name"`
	LastName string `json:"last_name" form:"last_name"`
	Active bool `json:"-"`
}

func NewUser(email string) *User {
	return &User{Email: email}
}

func GetUser(email string) *User {
	var user User
	DBConn.Find(&user, "email = ?", &email)
	return &user
}

func (u *User) Save() {
	db:=DBConn
	fmt.Println(db)
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(u)
}
