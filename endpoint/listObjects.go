package endpoint

import (
	"aws-s3-server/contracts"
	"aws-s3-server/types"
	"encoding/xml"
	"net/http"
	"strconv"
)

func ListObjects(DB contracts.Database, w http.ResponseWriter, r *http.Request, path string) {

	list, check := validatePathAndRequestForListObjects(path, r, w)
	if check == false {
		return
	}

	list.Contents, _ = DB.SelectContents(list.Name, list.Prefix)

	xmlBytes, err := xml.MarshalIndent(list, "", "  ")
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

func validatePathAndRequestForListObjects(path string, r *http.Request, w http.ResponseWriter) (*types.ListBucketResult, bool) {
	bucketName := path[len("/"):]
	if len(bucketName) == 0 {
		http.Error(w, "Wrong bucket", http.StatusUnprocessableEntity)
		return nil, false
	}
	query := r.URL.Query()

	marker := ""
	prefix := ""
	maxKeys := int64(0)

	markerValue := query["marker"]
	prefixValue := query["prefix"]
	maxKeysValue := query["max-keys"]

	if len(markerValue) > 0 {
		marker = markerValue[0]
	}

	if len(prefixValue) > 0 {
		prefix = prefixValue[0]
	}

	if len(maxKeysValue) > 0 {
		maxKeys, _ = strconv.ParseInt(maxKeysValue[0], 10, 12)
	}

	return &types.ListBucketResult{
		Name:        bucketName,
		IsTruncated: false,
		Marker:      marker,
		MaxKeys:     maxKeys,
		Prefix:      prefix,
	}, true
}
