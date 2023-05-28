package contracts

import "aws-s3-server/types"

type Database interface {
	BucketExists(name string) (int, error)
	InsertBucket(bucket *types.Bucket) error
	SelectContents(bucketName string, prefix string) ([]types.Content, error)
	InsertOrUpdateObject(object *types.Object) error
	GetObject(bucketId int, pathKey string) (*types.Object, error)
}
