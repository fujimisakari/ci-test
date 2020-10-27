package s3

import (
	"context"
	"io"
	"path/filepath"
)

type fakeS3ClientImpl struct {
	GetObjectReaderFunc func(ctx context.Context, key string) (io.ReadCloser, error)
	CopyObjectFunc      func(ctx context.Context, key string, writer io.Writer) error
	PutObjectFunc       func(ctx context.Context, key string, reader io.Reader) error
	ListObjectKeysFunc  func(ctx context.Context, prefix string, ext string) ([]string, error)
}

func (c *fakeS3ClientImpl) GetObjectReader(ctx context.Context, key string) (io.ReadCloser, error) {
	return c.GetObjectReaderFunc(ctx, key)
}

func (c *fakeS3ClientImpl) CopyObject(ctx context.Context, key string, writer io.Writer) error {
	return c.CopyObjectFunc(ctx, key, writer)
}

func (c *fakeS3ClientImpl) PutObject(ctx context.Context, key string, reader io.Reader) error {
	return c.PutObjectFunc(ctx, key, reader)
}

func (c *fakeS3ClientImpl) ListObjectKeys(ctx context.Context, prefix string, ext string) ([]string, error) {
	keys, err := c.ListObjectKeysFunc(ctx, prefix, ext)
	if err != nil {
		return nil, err
	}

	// exclude unmatched file's extention.
	var v []string
	for _, key := range keys {
		if filepath.Ext(key) == ext {
			v = append(v, key)
		}
	}
	return v, nil
}
