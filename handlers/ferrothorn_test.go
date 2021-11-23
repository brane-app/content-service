package handlers

import (
	"io"
	"strings"
	"testing"
)

type NopReader struct{}

func (NopReader) Read(_ []byte) (_ int, err error) {
	err = io.EOF
	return
}

func Test_upload_prefix(test *testing.T) {
	backupHost := ferrothorn_host
	backupMask := ferrothorn_mask

	defer func() {
		ferrothorn_host = backupHost
		ferrothorn_mask = backupMask
	}()

	ferrothorn_host = "host"
	ferrothorn_mask = ""

	url, err := upload(NopReader{})

	if err != nil {
		test.Fatal(err)
	}

	if !strings.HasPrefix(url, ferrothorn_host) {
		test.Fatalf("url doesn't have %s prefix: %s", ferrothorn_mask, url)
	}
}

func Test_upload_prefixMasked(test *testing.T) {
	backupHost := ferrothorn_host
	backupMask := ferrothorn_mask

	defer func() {
		ferrothorn_host = backupHost
		ferrothorn_mask = backupMask
	}()

	ferrothorn_host = "host"
	ferrothorn_mask = "mask"

	url, err := upload(NopReader{})

	if err != nil {
		test.Fatal(err)
	}

	if strings.HasPrefix(url, ferrothorn_host) {
		test.Fatalf("url still has %s prefix: %s", ferrothorn_host, url)
	}

	if !strings.HasPrefix(url, ferrothorn_mask) {
		test.Fatalf("url doesn't have %s prefix: %s", ferrothorn_mask, url)
	}
}
