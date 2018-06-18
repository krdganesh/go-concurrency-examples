package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

// Posts : Its a array of struct Post. Which is used to unmarshal data of an static api
type Posts []Post

// Post : Its a struct. Which contains all the fields of a post.
type Post struct {
	UserID int64  `json:"userId"`
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// startAsyncWebServiceDemo : Function calls getPosts function and returns true at the end
// Input Parameters
//         ch chan bool : boolen channel
func startAsyncWebServiceDemo(ch chan bool) {
	startTime := time.Now()
	runtime.GOMAXPROCS(4)
	cPost := make(chan Post)
	go getPosts(cPost)

	for data := range cPost {
		fmt.Println(data.UserID, data.ID)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nExcecution time: %s", elapsed)

	ch <- true
}

// getPosts : Function gets all posts from an external api and calls savePost function for each post recieved.
// Input Parameters
//         cPost chan Post : Post channel
func getPosts(cPost chan Post) {
	resp, _ := http.Get("https://jsonplaceholder.typicode.com/posts")
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var posts Posts
	json.Unmarshal(body, &posts)
	for _, post := range posts {
		if post.Body != "" {
			c := make(chan Post)
			go post.savePost(c)
			data := <-c
			cPost <- data
		} else {
			fmt.Println("Fatal Error - Data not found in body " + strconv.FormatInt(post.UserID, 10) + "_" + strconv.FormatInt(post.ID, 10))
		}
	}
	close(cPost)
}

// savePost : Function saves the post (json) in a file.
// Input Parameters
//         cPost chan Post : Post channel
func (post Post) savePost(c chan Post) {
	var fileName string
	fileName = strconv.FormatInt(post.UserID, 10) + "_" + strconv.FormatInt(post.ID, 10) + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	f, err := os.Create(watchedPath + "/" + fileName)
	if err != nil {
		fmt.Println(err)
	}
	jsonData, _ := json.Marshal(post)
	n, err := f.WriteString(string(jsonData))
	if err != nil {
		fmt.Println(err, n)
	}
	f.Close()
	c <- post
}
