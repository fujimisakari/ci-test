package s3

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client interface {
	GetObjectReader(ctx context.Context, key string) (io.ReadCloser, error)
	CopyObject(ctx context.Context, key string, writer io.Writer) error
	PutObject(ctx context.Context, key string, reader io.Reader) error
	ListObjectKeys(ctx context.Context, prefix string, ext string) ([]string, error)
}

var (
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Endpoint        string
	Bucket          string
)

// NewClient returns a Client.
func NewClient(isFake bool) (Client, error) {
	if isFake {
		return NewFakeClient(), nil
	}

	creds := credentials.NewStaticCredentials(AccessKeyID, SecretAccessKey, "")
	awsConf := aws.NewConfig()
	awsConf.WithCredentials(creds)
	awsConf.WithRegion(Region)
	// endpoint is determined automatically by aws-sdk-go
	// but local testing uses localhost.
	if len(Endpoint) != 0 {
		awsConf.WithEndpoint(Endpoint)
		awsConf.WithDisableSSL(true)
		awsConf.WithS3ForcePathStyle(true)
	}

	sess, err := session.NewSession(awsConf)
	if err != nil {
		return nil, err
	}
	return &s3ClientImpl{S3: s3.New(sess), bucket: Bucket}, nil
}

// NewFakeClient returns a new Client of Fake.
func NewFakeClient() *fakeS3ClientImpl { return &fakeS3ClientImpl{} }
