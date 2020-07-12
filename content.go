package main

import (
	"github.com/gastrodon/groudon"
	"github.com/imonke/monkebase"
	"github.com/imonke/monketype"

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
	ferrothorn_secret             = os.Getenv("FERROTHORN_SECRET")
	requester         http.Client = http.Client{}
)

func ferroRequest(sendable *http.Request) (response map[string]string, err error) {
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
		url = ferrothorn_host + "/" + response["id"]
	}

	return
}

func multipartReader(form *multipart.Form, key string) (reader io.Reader, ok bool) {
	var values []string
	if values, ok = form.Value[key]; !ok || len(values) != 1 {
		return
	}

	reader = strings.NewReader(values[0])
	return
}

func makeContent(data, file io.Reader, author string) (created monketype.Content, ok bool, err error) {
	var body CreateContentBody
	var external error
	if err, external = groudon.SerializeBody(data, &body); err != nil || external != nil {
		ok = external == nil
		return
	}

	var file_url string
	if file_url, err = upload(file); err != nil {
		return
	}

	created = monketype.NewContent(
		file_url,
		author,
		body.Mime,
		body.Tags,
		body.Featurable,
		body.NSFW,
	)

	err = monkebase.WriteContent(created.Map())
	ok = true
	return
}

func postContent(request *http.Request) (code int, r_map map[string]interface{}, err error) {
	var data io.Reader
	var ok bool
	if data, ok = multipartReader(request.MultipartForm, "json"); !ok {
		code = 400
		return
	}

	var file multipart.File
	var header *multipart.FileHeader
	if file, header, err = request.FormFile("file"); err != nil || header == nil {
		err = nil
		code = 400
		return
	}

	var created monketype.Content
	var author string = request.Context().Value("requester").(string)
	created, ok, err = makeContent(data, file, author)

	switch {
	case err != nil:
	case !ok:
		code = 400
	default:
		code = 200
		r_map = map[string]interface{}{"content": created}
	}

	return
}