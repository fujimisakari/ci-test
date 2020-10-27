package downloader

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fujimisakari/ci-test/s3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type downloader struct {
	s3Client  s3.Client
	keyPrefix string
	extension string
	logger    *zap.Logger
}

func NewDownloader(s3Client s3.Client, keyPrefix, extension string, logger *zap.Logger) *downloader {
	return &downloader{
		s3Client:  s3Client,
		keyPrefix: keyPrefix,
		extension: extension,
		logger:    logger,
	}
}

func (d *downloader) Do(ctx context.Context) ([]string, error) {
	objKeys, err := d.s3Client.ListObjectKeys(ctx, d.keyPrefix, d.extension)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get S3 object keys")
	}

	errs := d.download(ctx, objKeys)
	if len(errs) > 0 {
		for _, e := range errs {
			err = errors.Wrap(e, "failed to download To local file")
		}
		return objKeys, err
	}
	return objKeys, nil
}

func (d *downloader) download(ctx context.Context, objKeys []string) []error {
	var wg sync.WaitGroup
	errch := make(chan error)
	sem := make(chan struct{}, 10)

	for _, key := range objKeys {
		wg.Add(1)

		go func(k string) {
			sem <- struct{}{}
			defer func() { <-sem }()
			defer wg.Done()

			d.logger.Info(fmt.Sprintf("download from S3: %s", k))

			f, err := os.Create(FilePathFromObjKey(k))
			if err != nil {
				errch <- err
				return
			}
			defer f.Close()

			err = d.s3Client.CopyObject(ctx, k, f)
			if err != nil {
				errch <- err
				return
			}
		}(key)
	}

	go func() {
		wg.Wait()
		close(errch)
	}()

	var errs []error
	for err := range errch {
		errs = append(errs, err)
	}

	return errs
}

func FilePathFromObjKey(key string) string {
	return fmt.Sprintf("/tmp/%s", strings.Replace(key, "/", "_", -1))
}
