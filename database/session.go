package database

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model
	Token string
	UserId uint
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (s *Session)Valid() bool {
	if s.ExpiresAt.Sub(s.CreatedAt) < 0 {
		return false
	}
	if time.Now().Sub(s.CreatedAt) < 0 {
		return false
	}
	x:=time.Now().Sub(s.ExpiresAt)
	fmt.Println(x)
	if time.Now().Sub(s.ExpiresAt) > 0 {
		return false
	}
	return true
}