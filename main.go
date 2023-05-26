package main

import (
	"aws-s3-server/database"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type ListBucketResult struct {
	XMLName     xml.Name  `xml:"ListBucketResult"`
	Name        string    `xml:"Name"`
	Prefix      string    `xml:"Prefix"`
	Marker      string    `xml:"Marker"`
	MaxKeys     int       `xml:"MaxKeys"`
	IsTruncated bool      `xml:"IsTruncated"`
	Contents    []Content `xml:"Contents"`
}

type Content struct {
	XMLName      xml.Name  `xml:"Content"`
	Key          string    `xml:"Key"`
	LastModified time.Time `xml:"LastModified"`
	ETag         string    `xml:"ETag"`
	Size         int       `xml:"Size"`
}

var buckets = make(map[string]ListBucketResult)

var DB *database.Database

func main() {
	DB = database.Builder()
	http.HandleFunc("/", CreateBucket)
	http.HandleFunc("/", PutObject)
	http.HandleFunc("/", ListObjects)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if r.Method != http.MethodPut || strings.Count(path, "/") != 1 {
		return
	}

	bucketName := path[len("/"):]
	var bucket = &database.Bucket{
		Name: bucketName,
	}

	// TODO check if the bucket just exists

	err1 := DB.InsertBucket(bucket)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err2 := fmt.Fprintf(w, "Bucket %s created\n", bucketName)
	if err2 != nil {
		return
	}
}

func PutObject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.URL.Path[len("/Bucket/"):]
	_, ok := buckets[bucketName]
	if !ok {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return
	}

	// Here, the implementation of the object addition to the bucket should occur.
	// This part is missing as the handling of the file upload, its validation,
	// and its subsequent storage is dependent on your specific use case and setup.

	w.WriteHeader(http.StatusOK)
	// Assuming ETag is a computed value based on the object.
	_, err := fmt.Fprintf(w, `ETag: "%s"`, "your-etag-value")
	if err != nil {
		return
	}
}

func ListObjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.URL.Path[len("/Bucket/"):]
	bucket, ok := buckets[bucketName]
	if !ok {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return
	}

	// List objects from bucket. In this case, just marshalling the bucket object into XML
	xmlBytes, err := xml.MarshalIndent(bucket, "", "  ")
	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	_, write := w.Write(xmlBytes)
	if write != nil {
		return
	}
}
