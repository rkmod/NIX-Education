package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	postsCount = 100

	urlPostsScheme = "http"
	urlPostsHost   = "jsonplaceholder.typicode.com"
	urlPostsPath   = "posts/"
)

// reaches the content on the website, records it to the memory and prints it in the terminal.
func f(i int, postsURL string) {

	resp, err := http.Get(postsURL)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	posts := string(body)
	fmt.Printf(posts)
}

// iteratates through the links and launches goroutines
func main() {

	for i := 0; i < postsCount; i++ {
		postNumber := strconv.Itoa(i + 1)

		postsURL := url.URL{
			Scheme: urlPostsScheme,
			User:   &url.Userinfo{},
			Host:   urlPostsHost,
			Path:   urlPostsPath + postNumber,
		}
		amt := time.Duration(rand.Intn(2500))
		time.Sleep(time.Microsecond * amt)
		go f(i, postsURL.String())
	}

	input := ""
	fmt.Scanln(&input)
}
