package endpoint

import (
	"aws-s3-server/contracts"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func validateBucket(DB contracts.Database, w http.ResponseWriter, bucketName string) (int, bool) {
	bucketId, _ := DB.BucketExists(bucketName)
	if bucketId == 0 {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return 0, false
	}

	return bucketId, true
}

func dataFn(r *http.Request) (string, []byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Body not closed")
		}
	}(r.Body)

	hash := md5.Sum(body)

	base64Hash := base64.StdEncoding.EncodeToString(hash[:])

	return base64Hash, body, nil
}

func validatePathWhenBucketAndKeyIsPresent(w http.ResponseWriter, path string) ([]string, bool) {
	pathString := path[len("/"):]
	pathData := strings.Split(pathString, "/")
	if len(pathData) < 2 || len(pathData[0]) == 0 {
		http.Error(w, fmt.Sprintf("Bucket or path %s is invalid \n", pathString), http.StatusUnprocessableEntity)
		return nil, false
	}
	bucket := pathData[0]
	pathData = pathData[1:]

	for _, value := range pathData {
		if len(value) == 0 {
			http.Error(w, fmt.Sprintf("Key or path %s is invalid \n", pathString), http.StatusUnprocessableEntity)
			return nil, false
		}
	}

	data := []string{bucket, strings.Join(pathData, "/")}
	return data, true
}
