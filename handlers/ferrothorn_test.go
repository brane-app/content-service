package handlers

import (
	"io"
	"testing"
)

type NopReader struct{}

func (NopReader) Read(_ []byte) (_ int, err error) {
	err = io.EOF
	return
}

func Test_upload_prefix(test *testing.T) {
	if _, err := upload(NopReader{}); err != nil {
		test.Fatal(err)
	}
}
