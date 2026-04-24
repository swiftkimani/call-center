package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/yourorg/callcenter/internal/config"
)

type Client struct {
	mc     *minio.Client
	bucket string
}

func New(cfg *config.Config) (*Client, error) {
	mc, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio.New: %w", err)
	}
	return &Client{mc: mc, bucket: cfg.MinioBucket}, nil
}

func (c *Client) EnsureBucket(ctx context.Context) error {
	exists, err := c.mc.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("BucketExists: %w", err)
	}
	if !exists {
		if err := c.mc.MakeBucket(ctx, c.bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("MakeBucket: %w", err)
		}
	}
	return nil
}

func (c *Client) Upload(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	_, err := c.mc.PutObject(ctx, c.bucket, key, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("PutObject: %w", err)
	}
	return nil
}

func (c *Client) PresignedGetURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	u, err := c.mc.PresignedGetObject(ctx, c.bucket, key, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("PresignedGetObject: %w", err)
	}
	return u.String(), nil
}
