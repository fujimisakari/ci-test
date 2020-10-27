package downloader

import (
	"context"
	"io"
	"testing"

	"github.com/fujimisakari/ci-test/s3"
	"go.uber.org/zap"
)

func TestDo(t *testing.T) {
	fakeS3Client := s3.NewFakeClient()
	fakeS3Client.ListObjectKeysFunc = func(ctx context.Context, prefix string, ext string) ([]string, error) {
		return []string{
			"testdata/2019/",
			"testdata/2019/01/",
			"testdata/2019/01/10/",
			"testdata/2019/01/10/1.txt",
			"testdata/2019/01/10/2.csv",
			"testdata/2019/01/10/3.txt",
		}, nil
	}
	fakeS3Client.CopyObjectFunc = func(ctx context.Context, key string, writer io.Writer) error {
		return nil
	}

	downloader := NewDownloader(fakeS3Client, "testdata/2019/01/10/", ".txt", zap.NewNop())
	objKeys, err := downloader.Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(objKeys), 2; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := objKeys[0], "testdata/2019/01/10/1.txt"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
	if got, want := objKeys[1], "testdata/2019/01/10/3.txt"; got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
