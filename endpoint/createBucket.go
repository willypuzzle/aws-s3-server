package endpoint

import (
	"aws-s3-server/database"
	"aws-s3-server/types"
	"fmt"
	"net/http"
)

func CreateBucket(DB *database.Database, w http.ResponseWriter, path string) {
	bucketName := path[len("/"):]
	check, _ := DB.BucketExists(bucketName)

	if check != 0 {
		http.Error(w, fmt.Sprintf("Bucket %s just exists\n", bucketName), http.StatusUnprocessableEntity)
		return
	}

	var bucket = &types.Bucket{
		Name: bucketName,
	}

	err1 := DB.InsertBucket(bucket)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(fmt.Sprintf("Bucket %s created\n", bucketName))
}
