package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/minio/minio-go/v7"
)

// ObjectMeta is a small subset of object metadata for HTTP responses.
type ObjectMeta struct {
	ContentType string
	Size        int64
}

// Bucket wraps a MinIO client for a single bucket name.
type Bucket struct {
	client *minio.Client
	name   string
}

func NewBucket(client *minio.Client, bucket string) *Bucket {
	return &Bucket{client: client, name: bucket}
}

func (b *Bucket) EnsureExists(ctx context.Context) error {
	if b == nil || b.client == nil {
		return fmt.Errorf("storage: nil client")
	}
	ok, err := b.client.BucketExists(ctx, b.name)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	if err := b.client.MakeBucket(ctx, b.name, minio.MakeBucketOptions{}); err != nil {
		return err
	}
	slog.Info("minio: created bucket", "bucket", b.name)
	return nil
}

func (b *Bucket) Put(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) error {
	if b == nil || b.client == nil {
		return fmt.Errorf("storage: nil client")
	}
	opts := minio.PutObjectOptions{}
	if contentType != "" {
		opts.ContentType = contentType
	} else {
		opts.ContentType = "application/octet-stream"
	}
	_, err := b.client.PutObject(ctx, b.name, objectKey, reader, size, opts)
	return err
}

func (b *Bucket) Remove(ctx context.Context, objectKey string) error {
	if b == nil || b.client == nil || objectKey == "" {
		return nil
	}
	return b.client.RemoveObject(ctx, b.name, objectKey, minio.RemoveObjectOptions{})
}

// Open streams an object from the bucket. Caller must close the returned ReadCloser.
func (b *Bucket) Open(ctx context.Context, objectKey string) (io.ReadCloser, ObjectMeta, error) {
	var zero ObjectMeta
	if b == nil || b.client == nil || objectKey == "" {
		return nil, zero, fmt.Errorf("storage: invalid bucket or key")
	}
	obj, err := b.client.GetObject(ctx, b.name, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, zero, err
	}
	info, err := obj.Stat()
	if err != nil {
		_ = obj.Close()
		return nil, zero, err
	}
	return obj, ObjectMeta{ContentType: info.ContentType, Size: info.Size}, nil
}

// BucketName returns the configured bucket name.
func (b *Bucket) BucketName() string {
	if b == nil {
		return ""
	}
	return b.name
}
