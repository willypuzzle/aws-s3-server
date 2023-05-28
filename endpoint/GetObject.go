package endpoint

import (
	"aws-s3-server/contracts"
	"aws-s3-server/types"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

func GetObject(DB contracts.Database, w http.ResponseWriter, r *http.Request, path string) {
	pathData, check1 := validatePathWhenBucketAndKeyIsPresent(w, path)
	if check1 == false {
		return
	}

	bucketId, check2 := validateBucket(DB, w, pathData[0])
	if check2 == false {
		return
	}

	keyPath := pathData[1]

	data, rangeValue, check3 := validateHeadersAndGetDataForGetObject(DB, bucketId, keyPath, r, w)
	if check3 == false {
		return
	}

	w.Header().Set("Last-Modified", data.UpdatedAt.String())
	w.Header().Set("ETag", data.Uuid)
	w.Header().Set("Content-Range", fmt.Sprintf("bytes=%d-%d", rangeValue[0], rangeValue[1]))
	w.Header().Set("Content-Type", data.ContentType)
	_, write := w.Write(data.Data)
	if write != nil {
		return
	}
}

func validateHeadersAndGetDataForGetObject(
	DB contracts.Database,
	bucketId int,
	keyPath string,
	r *http.Request,
	w http.ResponseWriter,
) (*types.Object, [2]int, bool) {

	data, err1 := DB.GetObject(bucketId, keyPath)
	if err1 != nil {
		http.Error(w, "Recovery Object Error", http.StatusInternalServerError)
		return nil, [2]int{0, 0}, false
	}
	if data == nil {
		http.Error(w, "Object not found", http.StatusNotFound)
		return nil, [2]int{0, 0}, false
	}
	var rangeArray [2]int
	rangeArray, data.Data = parseRangeHeaderAndCustomizeData(data.Data, r)

	return data, rangeArray, true
}

func parseRangeHeaderAndCustomizeData(data []byte, r *http.Request) ([2]int, []byte) {
	var rangeArray [2]int
	rangeValueHeader := r.Header["Range"]
	if rangeValueHeader == nil || len(rangeValueHeader) < 1 {
		rangeArray[0] = 0
		rangeArray[1] = len(data)
	} else {
		var check bool
		rangeArray, check = parseRangeHeaderString(rangeValueHeader[0])
		if check == false {
			rangeArray[0] = 0
			rangeArray[1] = len(data)
		}
	}

	return rangeArray, data[rangeArray[0]:rangeArray[1]]
}

func parseRangeHeaderString(s string) ([2]int, bool) {
	r := regexp.MustCompile(`bytes=(\d+)-(\d+)`)
	matches := r.FindStringSubmatch(s)
	if len(matches) != 3 {
		return [2]int{0, 0}, false
	}
	a, err1 := strconv.ParseInt(matches[1], 10, 12)
	if err1 != nil {
		return [2]int{0, 0}, false
	}
	b, err2 := strconv.ParseInt(matches[2], 10, 12)
	if err2 != nil {
		return [2]int{0, 0}, false
	}
	return [2]int{int(a), int(b)}, true
}
