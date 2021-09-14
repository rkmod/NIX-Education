package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

// Comment defines all fields related to comment
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

type dbConnection struct {
	DB    *gorm.DB
	sqlDB *gorm.DB
	Close *gorm.DB
}

// NewConnection opens connection to DB and returns it
func NewDBConnection() (*dbConnection, error) {
	dsn := "***:***@tcp(127.0.0.1:3306)/comments"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect to the database")
	}

	return &dbConnection{DB: db}, nil
}

// records posts into the database
func (sqlDB *dbConnection) setPostIntoDB(post Post) error {
	posts := Post{UserId: post.UserId, Title: post.Title, Body: post.Body}
	sqlDB.DB.Create(&posts)

	return nil
}

// records comments into the database
func (sqlDB *dbConnection) setCommentIntoDB(comment Comment) error {
	comments := Comment{PostId: comment.PostId, Id: comment.Id, Name: comment.Name, Email: comment.Email, Body: comment.Body}
	sqlDB.DB.Create(&comments)

	return nil
}

// reaches the website with the posts, iterates through the websites with comments and launches goroutines
func main() {
	var wg sync.WaitGroup
	sqlDB, err := NewDBConnection()
	if err != nil {
		log.Fatalf("failed to connection to DB: error: %s", err.Error())
	}

	dbURL := url.URL{
		Scheme:   "https",
		Host:     "jsonplaceholder.typicode.com",
		Path:     "posts",
		RawQuery: "userId=7",
	}

	resp, err := http.Get(dbURL.String())
	if err != nil {
		log.Fatalln(err)
	}

	commentBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	posts := make([]Post, 0)
	if err := json.Unmarshal(commentBody, &posts); err != nil {
		panic(err)
	}

	for id, post := range posts {
		err := sqlDB.setPostIntoDB(post)
		if err != nil {
			log.Fatalf("failed to set post into DB: %v", err)
		}

		fmt.Printf("The last inserted row id: %d\n", id)
	}

	dbURL.Path = "comments"

	for i := 1; i < 101; i++ {
		dbURL.RawQuery = "postId=" + strconv.Itoa(i)

		wg.Add(1)
		go func() {
			defer wg.Done()

			b := make(chan []byte)

			go getComments(b, dbURL.String())
			go writeCommentsIntoDb(sqlDB, b)

			fmt.Println("Finished")
		}()

		time.Sleep(time.Microsecond * time.Duration(rand.Intn(50000)))
	}

	wg.Wait()
	input := ""
	fmt.Scanln(&input)
}

func getComments(b chan []byte, url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	commentBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Finished getComments")

	b <- commentBody
}

// records comments into the database
func writeCommentsIntoDb(sqlDB *dbConnection, b chan []byte) {
	receivedCommentBody := <-b

	comments := make([]Comment, 0)
	if err := json.Unmarshal(receivedCommentBody, &comments); err != nil {
		panic(err)
	}

	for id, comment := range comments {
		err := sqlDB.setCommentIntoDB(comment)
		if err != nil {
			log.Fatalf("failed to set comment into DB: %v", err)
		}

		fmt.Printf("The last inserted row id: %d\n", id)
	}

	fmt.Println("Finished writeComments")
}
