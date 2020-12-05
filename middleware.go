package main

import (
	"github.com/imonke/monkelib/middleware"

	"bytes"
	"encoding/base64"
	"io"
	"mime/multipart"
	"net/http"
)

func transformBase64(request *http.Request) (modified *http.Request, ok bool, _ int, r_map map[string]interface{}, err error) {
	var values []string
	var exists bool
	if values, exists = request.MultipartForm.Value["file_base64"]; !exists || len(values) == 0 {
		ok = true
		return
	}

	var data_bytes []byte
	if data_bytes, err = base64.StdEncoding.DecodeString(values[0]); err != nil {
		err = nil
		return
	}

	var body *bytes.Buffer = new(bytes.Buffer)
	var writer *multipart.Writer = multipart.NewWriter(body)

	var part_file io.Writer
	if part_file, err = writer.CreateFormFile("file", "file_base64_transformed"); err != nil {
		return
	}
	part_file.Write(data_bytes)

	if values, exists = request.MultipartForm.Value["json"]; exists || len(values) != 0 {
		writer.WriteField("json", values[0])
	}

	writer.Close()

	if modified, err = http.NewRequestWithContext(request.Context(), request.Method, request.URL.String(), body); err != nil {
		err = nil
		return
	}
	modified.Header.Set("Content-Type", writer.FormDataContentType())
	err = modified.ParseMultipartForm(middleware.MULTIPART_MEM_MAX)
	ok = true
	return
}
