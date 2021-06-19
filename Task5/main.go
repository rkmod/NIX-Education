package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"
)

const (
	postsPath           = "/Users/Asus/go/src/.github.io/NIX-Education/Task5/storage/posts/"
	postsCount          = 100
	fileExstension      = ".txt"
	filePermissionsCode = 0666

	fileRecorded = "File %s recorded\n"
	finishedAll  = "Finished All. Please press Ctrl + C to exit the program."

	urlPostsScheme = "http"
	urlPostsHost   = "jsonplaceholder.typicode.com"
	urlPostsPath   = "posts/"
)

// Post defines all fields related to post
type Post struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// reaches the content on the website and records it in the memory
func getPostInfo(c chan []byte, postsURL string) {
	resp, err := http.Get(postsURL)
	if err != nil {
		log.Fatalf("failed to get info by posts url: %v", err)
	}

	postBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body with posts: %v", err)
	}

	c <- postBody

}

// passes the recorded content and writes it to the file
func writePostIntoFile(c chan []byte, definitePath string) {
	receivedJsPost := <-c

	var post Post

	if err := json.Unmarshal(receivedJsPost, &post); err != nil {
		panic(err)
	}

	postStr := fmt.Sprintf("%v", post) // makes it possible to record values of the struct without keys to the file.

	/*In order to record the struct with keys into the file, please comment line 66,
	write this line of the code next: mPost, _ := json.Marshal(post)
	and replace []byte(postStr with mPost in the line 75.*/

	file, err := os.OpenFile(definitePath, os.O_RDWR, filePermissionsCode)
	if err != nil {
		log.Fatalf("failed to open file %s, error: %v", definitePath, err)
	}

	err = ioutil.WriteFile(definitePath, []byte(postStr), filePermissionsCode)
	if err != nil {
		log.Fatalf("failed to write into a file %s, error: %v", definitePath, err)
	}

	file.Close()
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < postsCount; i++ {
		postNumber := strconv.Itoa(i + 1)
		fileName := postNumber + fileExstension
		definitePath := path.Join(postsPath, fileName)

		_, err := os.Stat(definitePath)
		if os.IsNotExist(err) {
			file, err := os.Create(definitePath)
			if err != nil {
				log.Fatalf("failed to create file %s, error: %v", definitePath, err)
			}
			file.Close()
		}

		fmt.Printf(fileRecorded, fileName) // message that acknowledges that the particular file has been recorded

		postsURL := url.URL{
			Scheme: urlPostsScheme,
			User:   &url.Userinfo{},
			Host:   urlPostsHost,
			Path:   urlPostsPath + postNumber,
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			c := make(chan []byte)

			go getPostInfo(c, postsURL.String())
			go writePostIntoFile(c, definitePath)
		}()
	}

	wg.Wait()

	fmt.Println(finishedAll)

	input := ""
	fmt.Scanln(&input)
}
