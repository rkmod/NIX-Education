package main

import (
	"database/sql"
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

	_ "github.com/go-sql-driver/mysql"
)

// Post defines all fields related to post
type Post struct {
	Num    int    `json:"" db:"Num"`
	UserId int    `json:"userId" db:"user_Id"`
	Title  string `json:"title" db:"title"`
	Body   string `json:"body" db:"body"`
}

// Comment defines all fields related to comment
type Comment struct {
	Num    int    `db:"Num"`
	PostId int    `json:"postId" db:"post_Id"`
	Id     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Email  string `json:"email" db:"email"`
	Body   string `json:"body" db:"body"`
}

type dbConnection struct {
	DB *sql.DB
}

// NewConnection opens connection to DB and returns it
func NewDBConnection() (*dbConnection, error) {
	db, err := sql.Open("mysql", "nameOfTheDB:password@tcp(127.0.0.1:3306)/comments")
	if err != nil {
		return nil, err
	}

	return &dbConnection{DB: db}, nil
}

// records posts into the database
func (db *dbConnection) setPostIntoDB(post Post) error {
	stmt, err := db.DB.Prepare("insert into posts(user_Id, title, body) values (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(post.UserId, post.Title, post.Body)
	if err != nil {
		return err
	}

	return nil
}

// records comments into the database
func (db *dbConnection) setCommentIntoDB(comment Comment) error {
	stmt, err := db.DB.Prepare("insert into comments(post_Id, id, name, email, body) values (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(comment.PostId, comment.Id, comment.Name, comment.Email, comment.Body)
	if err != nil {
		return err
	}

	return nil
}

// reaches the website with the posts, iterates through the websites with comments and launches goroutines
func main() {
	var wg sync.WaitGroup

	db, err := NewDBConnection()
	if err != nil {
		log.Fatalf("failed to connection to DB: error: %s", err.Error())
	}
	defer db.DB.Close()

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
		err := db.setPostIntoDB(post)
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
			go writeCommentsIntoDb(db, b)

			fmt.Println("Finished")
		}()

		time.Sleep(time.Microsecond * time.Duration(rand.Intn(50000)))
	}

	wg.Wait()

	input := ""
	fmt.Scanln(&input)
}

// reaches the website, gets coomments and records them into the memory
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
func writeCommentsIntoDb(db *dbConnection, b chan []byte) {
	receivedCommentBody := <-b

	comments := make([]Comment, 0)
	if err := json.Unmarshal(receivedCommentBody, &comments); err != nil {
		panic(err)
	}

	for id, comment := range comments {
		err := db.setCommentIntoDB(comment)
		if err != nil {
			log.Fatalf("failed to set comment into DB: %v", err)
		}

		fmt.Printf("The last inserted row id: %d\n", id)
	}

	fmt.Println("Finished writeComments")
}
