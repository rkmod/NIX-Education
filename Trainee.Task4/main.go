package main

import (
	"log"
	"net/http"

	"Trainee.Task4/controllers"
	"Trainee.Task4/mydatabase"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	initDB()
	log.Println("Starting the HTTP server on port 8090")

	router := mux.NewRouter().StrictSlash(true)
	initaliseHandlersForPosts(router)
	initaliseHandlersForComments(router)
	log.Fatal(http.ListenAndServe(":8090", router))
}

func initaliseHandlersForPosts(router *mux.Router) {
	router.HandleFunc("/create/post", controllers.CreatePost).Methods("POST")
	router.HandleFunc("/get/post", controllers.GetAllPost).Methods("GET")
	router.HandleFunc("/get/post/{num}", controllers.GetPostByID).Methods("GET")
	router.HandleFunc("/update/post/{num}", controllers.UpdatePostByID).Methods("PUT")
	router.HandleFunc("/delete/post/{num}", controllers.DeletePostByID).Methods("DELETE")
}

func initaliseHandlersForComments(router *mux.Router) {
	router.HandleFunc("/create/comment", controllers.CreateComment).Methods("POST")
	router.HandleFunc("/get/comment", controllers.GetAllComment).Methods("GET")
	router.HandleFunc("/get/comment/{num}", controllers.GetCommentByID).Methods("GET")
	router.HandleFunc("/update/comment/{num}", controllers.UpdateCommentByID).Methods("PUT")
	router.HandleFunc("/delete/comment/{num}", controllers.DeleteCommentByID).Methods("DELETE")
}

func initDB() {
	config :=
		mydatabase.Config{
			ServerName: "127.0.0.1:3306",
			User:       "***",
			Password:   "***",
			DB:         "comments",
		}

	connectionString := mydatabase.GetConnectionString(config)
	err := mydatabase.NewDBConnection(connectionString)
	if err != nil {
		panic(err.Error())
	}
}
