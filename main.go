package main

import (
	"aws-s3-server/database"
	"aws-s3-server/endpoint"
	"fmt"
	"net/http"
	"strings"
)

var DB *database.Database

func main() {
	DB = database.Builder()
	http.HandleFunc("/", manageRoute)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server Hung up")
		return
	}
}

func manageRoute(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(fmt.Sprintf("Request received on %s\n", path))
	if r.Method == http.MethodPut && strings.Count(path, "/") == 1 {
		endpoint.CreateBucket(DB, w, path)
	} else if r.Method == http.MethodPut && strings.Count(path, "/") >= 2 {
		endpoint.PutObject(DB, w, r, path)
	} else if r.Method == http.MethodGet && strings.Count(path, "/") == 1 {
		endpoint.ListObjects(DB, w, r, path)
	} else {
		http.Error(w, "Unknown API", http.StatusBadRequest)
	}
}
