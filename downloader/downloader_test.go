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

func TestFilePathFromObjKey(t *testing.T) {
	cases := []struct {
		objKey string
		want   string
	}{
		{
			"testdata/2018/12/01/1.txt",
			"/tmp/testdata_2018_12_01_1.txt",
		},
		{
			"testdata/2019/01/28/1.txt",
			"/tmp/testdata_2019_01_28_1.txt",
		},
		{
			"testdata/2020/02/28/1.txt",
			"/tmp/testdata_2020_02_28_1.txt",
		},
		{
			"testdata/2020/03/28/1.txt",
			"/tmp/testdata_2020_0_28_1.txt",
		},
	}

	for _, tc := range cases {
		if got, want := FilePathFromObjKey(tc.objKey), tc.want; got != want {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}
