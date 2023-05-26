package types

import "time"

type Bucket struct {
	Id   int
	Name string
}

type Object struct {
	Id          int
	BucketId    int
	Key         string
	Data        []byte
	Size        int
	ContentType string
	Uuid        string
	UpdatedAt   time.Time
}
