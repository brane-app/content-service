package main

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"testing"
)

const (
	pngBase64 = "iVBORw0KGgpoYWhhIHllcyBJIGFtIGEgcG5n"
)

func mustLocalBase64(data string) (request *http.Request) {
	var body *bytes.Buffer = new(bytes.Buffer)
	var writer *multipart.Writer = multipart.NewWriter(body)
	writer.WriteField("file_base64", data)
	writer.Close()

	var err error
	if request, err = http.NewRequest("GET", "http://localhost", body); err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.ParseMultipartForm(512)
	return
}

func Test_transformBase64(test *testing.T) {
	var request *http.Request
	var ok bool
	var err error
	if request, ok, _, _, err = transformBase64(mustLocalBase64(pngBase64)); err != nil {
		test.Fatal(err)
	}

	if !ok {
		test.Errorf("response for %s isn't ok", pngBase64)
	}

	var file multipart.File
	var header *multipart.FileHeader
	if file, header, err = request.FormFile("file"); err != nil || file == nil || header == nil {
		test.Errorf("file nil? %t\nheader nil? %t\n%v", file == nil, header == nil, err)
	}

	var file_bytes []byte
	if file_bytes, err = ioutil.ReadAll(file); err != nil {
		test.Fatal(err)
	}

	var decoded []byte
	decoded, _ = base64.StdEncoding.DecodeString(pngBase64)

	if string(file_bytes) != string(decoded) {
		test.Errorf("file data mismatch, expected %s, got %s", pngBase64, file_bytes)
	}
}
