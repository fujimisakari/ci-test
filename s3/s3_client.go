package s3

import (
	"context"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3ClientImpl struct {
	S3     *s3.S3
	bucket string
}

func (c *s3ClientImpl) GetObjectReader(ctx context.Context, key string) (io.ReadCloser, error) {
	v, err := c.S3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return v.Body, nil
}

func (c *s3ClientImpl) CopyObject(ctx context.Context, key string, writer io.Writer) error {
	r, err := c.GetObjectReader(ctx, key)
	if err != nil {
		return err
	}
	defer r.Close()
	_, err = io.Copy(writer, r)
	return err
}

func (c *s3ClientImpl) PutObject(ctx context.Context, key string, reader io.Reader) error {
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(reader),
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	_, err := c.S3.PutObjectWithContext(ctx, input)
	return err
}

func (c *s3ClientImpl) ListObjectKeys(ctx context.Context, prefix string, ext string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.bucket),
		Prefix: aws.String(prefix),
	}

	var ret []string

	for {
		output, err := c.S3.ListObjectsV2WithContext(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, obj := range output.Contents {
			key := aws.StringValue(obj.Key)
			if filepath.Ext(key) == ext {
				ret = append(ret, key)
			}
		}

		if output.NextContinuationToken == nil {
			break
		}

		input = &s3.ListObjectsV2Input{
			Bucket:            aws.String(c.bucket),
			Prefix:            aws.String(prefix),
			ContinuationToken: output.NextContinuationToken,
		}
	}

	return ret, nil
}
