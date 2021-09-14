package entity

import "time"

type Comment struct {
	Num       int       `json:"num" xml:"num" db:"Num" gorm:"primary_key`
	PostId    int       `json:"postId" xml:"postId" db:"post_Id"`
	Id        int       `json:"id" xml:"id" db:"id"`
	Name      string    `json:"name" xml:"name" db:"name"`
	Email     string    `json:"email" xml:"email" db:"email"`
	Body      string    `json:"body" xml:"body" db:"body"`
	CreatedAt time.Time `json:"created_at" xml:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" xml:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" xml:"deleted_at" db:"deleted_at"`
}
