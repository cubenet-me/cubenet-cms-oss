package main

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	client *minio.Client
	bucket string
}

type StorageConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func NewStorage(cfg StorageConfig) (*Storage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	ctx := context.Background()
	if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
		exists, errExists := client.BucketExists(ctx, cfg.Bucket)
		if errExists != nil || !exists {
			return nil, fmt.Errorf("create bucket: %w", err)
		}
	}

	return &Storage{client: client, bucket: cfg.Bucket}, nil
}

func (s *Storage) Upload(ctx context.Context, name string, r io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, s.bucket, name, r, size, minio.PutObjectOptions{})
	return err
}

func (s *Storage) Download(ctx context.Context, name string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *Storage) Delete(ctx context.Context, name string) error {
	return s.client.RemoveObject(ctx, s.bucket, name, minio.RemoveObjectOptions{})
}

func (s *Storage) URL(name string) string {
	return fmt.Sprintf("/api/v1/launcher/files/%s", name)
}
