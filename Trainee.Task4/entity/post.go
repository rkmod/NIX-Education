package entity

import (
	"time"
)

// Post defines all fields related to post
type Post struct {
	Num       int       `json:"num" xml:"num" db:"Num" gorm:"primary_key"`
	UserId    int       `json:"userId" xml:"userId" db:"user_Id"`
	Title     string    `json:"title" xml:"title" db:"title"`
	Body      string    `json:"body" xml:"body" db:"body"`
	CreatedAt time.Time `json:"created_at" xml:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" xml:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" xml:"deleted_at" db:"deleted_at"`
}
