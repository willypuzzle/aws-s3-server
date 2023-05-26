package endpoint

import (
	"aws-s3-server/database"
	"encoding/xml"
	"net/http"
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

func ListObjects(DB *database.Database, w http.ResponseWriter, r *http.Request, path string) {
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
