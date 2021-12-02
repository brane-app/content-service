package handlers

import (
	"github.com/google/uuid"

	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

var (
	ferrothorn_host               = strings.TrimSuffix(os.Getenv("FERROTHORN_HOST"), "/")
	ferrothorn_mask               = strings.TrimSuffix(os.Getenv("FERROTHORN_MASK"), "/")
	ferrothorn_secret             = os.Getenv("FERROTHORN_SECRET")
	noFerrothorn                  = os.Getenv("NO_FERROTHORN") == "true"
	requester         http.Client = http.Client{}
)

func ferroRequest(sendable *http.Request) (response map[string]string, err error) {
	if noFerrothorn {
		response = map[string]string{"id": "no_ferrothorn_" + uuid.New().String()}
		return
	}

	sendable.Header.Set("Authorization", ferrothorn_secret)

	var http_response *http.Response
	if http_response, err = requester.Do(sendable); err != nil {
		return
	}

	var data []byte
	if data, err = ioutil.ReadAll(http_response.Body); err != nil {
		return
	}

	if err = json.Unmarshal(data, &response); err == nil {
		var exists bool
		if _, exists = response["error"]; exists {
			err = errors.New(response["error"])
		}
	}

	return
}

func upload(file io.Reader) (url string, err error) {
	var closer io.Closer
	var ok bool
	if closer, ok = file.(io.Closer); ok {
		defer closer.Close()
	}

	var buffer bytes.Buffer
	var writer *multipart.Writer = multipart.NewWriter(&buffer)
	defer writer.Close()

	var writable io.Writer
	if writable, err = writer.CreateFormFile("file", "file"); err != nil {
		return
	}

	if _, err = io.Copy(writable, file); err != nil {
		return
	}

	writer.Close()
	var sendable *http.Request
	if sendable, err = http.NewRequest("POST", ferrothorn_host, &buffer); err != nil {
		return
	}

	sendable.Header.Set("Content-Type", writer.FormDataContentType())

	var response map[string]string
	if response, err = ferroRequest(sendable); err == nil {
		if ferrothorn_mask != "" {
			url = ferrothorn_mask + "/" + response["id"]
		} else {
			url = ferrothorn_host + "/" + response["id"]
		}
	}

	return
}
