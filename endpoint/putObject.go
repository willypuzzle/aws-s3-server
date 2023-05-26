package endpoint

import (
	"aws-s3-server/database"
	"aws-s3-server/types"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
)

func PutObject(DB *database.Database, w http.ResponseWriter, r *http.Request, path string) {
	data, check := validateAndBuildPutObject(DB, path, w, r)
	if check == false {
		return
	}

	err := DB.InsertOrUpdateObject(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error insert new object"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, `ETag: "%s"`, data.Uuid)
}

func validateAndBuildPutObject(DB *database.Database, path string, w http.ResponseWriter, r *http.Request) (*types.Object, bool) {
	pathData, check1 := validatePath(w, path)
	if check1 == false {
		return nil, false
	}

	data, contentType, check2 := validateHeaders(w, r)
	if check2 == false {
		return nil, false
	}

	bucketId, check3 := validateBucket(DB, w, pathData[0])
	if check3 == false {
		return nil, false
	}

	keyPath := pathData[1]
	object := data
	uuidValue := uuid.New()

	return &types.Object{
		BucketId:    bucketId,
		Key:         keyPath,
		Data:        object,
		ContentType: contentType,
		Uuid:        uuidValue.String(),
	}, true
}

func validateBucket(DB *database.Database, w http.ResponseWriter, bucketName string) (int, bool) {
	bucketId, _ := DB.BucketExists(bucketName)
	if bucketId == 0 {
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return 0, false
	}

	return bucketId, true
}

func validateHeaders(w http.ResponseWriter, r *http.Request) ([]byte, string, bool) {
	md5HeaderValueArray := r.Header["Content-Md5"]
	contentTypeHeaderValueArray := r.Header["Content-Type"]

	if len(md5HeaderValueArray) == 0 || len(md5HeaderValueArray[0]) == 0 {
		http.Error(w, fmt.Sprintf("Invalid Header(s) \n"), http.StatusUnprocessableEntity)
		return nil, "", false
	}

	var contentType string
	if len(contentTypeHeaderValueArray) == 0 || len(contentTypeHeaderValueArray[0]) == 0 {
		contentType = "application/octet-stream"
	} else {
		contentType = contentTypeHeaderValueArray[0]
	}

	md5Posted := md5HeaderValueArray[0]
	md5Calculated, data, hashsError := dataFn(r)
	if md5Calculated != md5Posted {
		http.Error(w, fmt.Sprintf("Uuid is invalid '%s' \n", md5Posted), http.StatusUnprocessableEntity)
		return nil, "", false
	}

	if hashsError != nil {
		http.Error(w, fmt.Sprintf("Hash error %s \n", hashsError.Error()), http.StatusInternalServerError)
		return nil, "", false
	}

	return data, contentType, true
}

func validatePath(w http.ResponseWriter, path string) ([]string, bool) {
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
