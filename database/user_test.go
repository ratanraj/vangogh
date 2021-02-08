package database

import (
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func initDB(url string) {
	DBConn, err := gorm.Open(sqlite.Open(url), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DBConn.AutoMigrate(&User{})
}

func TestUser(t *testing.T) {
	initDB("sqlite3.db")
	for _, tt := range []struct {
		Username string
		Password string
	}{
		{Username: "ratanraj", Password: "password123"},
		{Username: "mani", Password: "password123"},
		{Username: "admin", Password: "password222"},
	} {
		u := NewUser(tt.Username, tt.Password)
		u.Save()
	}
}
