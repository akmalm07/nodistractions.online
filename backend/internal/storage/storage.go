package gcs

import cloudstorage "cloud.google.com/go/storage"

// Storage is the boundary between file routes and the selected cloud SDK.
type Storage struct {
	Bucket *cloudstorage.BucketHandle
}

func NewStorage(client *cloudstorage.Client, bucket string) *Storage {
	return &Storage{Bucket: client.Bucket(bucket)}
}
