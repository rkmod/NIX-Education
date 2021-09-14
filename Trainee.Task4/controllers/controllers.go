package controllers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"Trainee.Task4/entity"
	"Trainee.Task4/mydatabase"

	"github.com/gorilla/mux"
)

//GetAllPost get all post data
func GetAllPost(w http.ResponseWriter, r *http.Request) {
	var posts []entity.Post
	mydatabase.Connector.Find(&posts)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
	xml.NewEncoder(w).Encode(posts)
}

//GetAllComments get all comment data
func GetAllComment(w http.ResponseWriter, r *http.Request) {
	var comments []entity.Comment
	mydatabase.Connector.Find(&comments)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
	xml.NewEncoder(w).Encode(comments)
}

//GetPostByID returns post with specific ID
func GetPostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["num"]

	var post entity.Post
	mydatabase.Connector.First(&post, key)
	fmt.Println("Key:", key)
	fmt.Println("Post", post)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
	xml.NewEncoder(w).Encode(post)
}

//GetCommentByID returns comment with specific ID
func GetCommentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["num"]

	var comment entity.Comment
	mydatabase.Connector.First(&comment, key)
	fmt.Println("Key:", key)
	fmt.Println("Comment", comment)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
	xml.NewEncoder(w).Encode(comment)
}

//CreatePost creates post
func CreatePost(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var post entity.Post
	json.Unmarshal(requestBody, &post)
	xml.Unmarshal(requestBody, &post)

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	post.DeletedAt = time.Now()

	mydatabase.Connector.Create(post)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
	xml.NewEncoder(w).Encode(post)
}

//CreateComment creates comment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var comment entity.Comment
	json.Unmarshal(requestBody, &comment)
	xml.Unmarshal(requestBody, &comment)

	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	comment.DeletedAt = time.Now()

	mydatabase.Connector.Create(comment)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
	xml.NewEncoder(w).Encode(comment)
}

//UpdatePostByID updates post with respective ID
func UpdatePostByID(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var post entity.Post
	json.Unmarshal(requestBody, &post)
	xml.Unmarshal(requestBody, &post)
	mydatabase.Connector.Save(&post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
	xml.NewEncoder(w).Encode(post)
}

//UpdateCommentByID updates comment with respective ID
func UpdateCommentByID(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var comment entity.Comment
	json.Unmarshal(requestBody, &comment)
	xml.Unmarshal(requestBody, &comment)
	mydatabase.Connector.Save(&comment)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
	xml.NewEncoder(w).Encode(comment)
}

//DeletPostByID delete's post with specific ID
func DeletePostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["num"]

	var post entity.Post
	num, _ := strconv.ParseInt(key, 10, 64)
	fmt.Println("num:", num)
	//mydatabase.Connector.Where("num = ?", num).Delete(&post) // SoftDelete
	mydatabase.Connector.Where("num = ?", num).Unscoped().Delete(&post) //Permanent deletion
	fmt.Println("comment", post)
	w.WriteHeader(http.StatusNoContent)
}

//DeletCommentByID delete's comment with specific ID
func DeleteCommentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["num"]

	var comment entity.Comment
	num, _ := strconv.ParseInt(key, 10, 64)
	fmt.Println("num:", num)
	//mydatabase.Connector.Where("num = ?", num).Delete(&comment) //SoftDelete
	mydatabase.Connector.Where("num = ?", num).Unscoped().Delete(&comment) //Permanent deleteion
	fmt.Println("comment", comment)
	w.WriteHeader(http.StatusNoContent)
}
