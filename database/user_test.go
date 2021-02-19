package database

import (
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func init() {
	initDB("sqlite3.db")
}

func initDB(url string) {
	var err error
	DBConn, err = gorm.Open(sqlite.Open(url), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DBConn.AutoMigrate(&User{})
}

func TestUser(t *testing.T) {
	for _, tt := range []struct {
		Username string
	}{
		{Username: "user01"},
		{Username: "admin"},
	} {
		u := NewUser(tt.Username)
		u.Save()
	}
}
