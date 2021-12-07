package handlers

import (
	"github.com/brane-app/content-service/types"

	"github.com/brane-app/database-library"
	library_types "github.com/brane-app/types-library"
	"github.com/gastrodon/groudon/v2"

	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

var (
	allowed_mime map[string]bool = map[string]bool{
		"image/png":  true,
		"image/jpg":  true,
		"image/jpeg": true,
		"image/webp": true,
	}
)

func multipartReader(form *multipart.Form, key string) (reader io.Reader, ok bool) {
	var values []string
	if values, ok = form.Value[key]; !ok || len(values) != 1 {
		return
	}

	reader = strings.NewReader(values[0])
	return
}

func makeContent(data, file io.Reader, author string) (created library_types.Content, ok bool, err error) {
	var body types.CreateContentBody
	var external error
	if err, external = groudon.SerializeBody(data, &body); err != nil || external != nil {
		ok = external == nil
		return
	}

	var file_tee bytes.Buffer
	file = io.TeeReader(file, &file_tee)

	var file_bytes []byte = make([]byte, 512)
	if _, err = file.Read(file_bytes); err != nil {
		return
	}

	ioutil.ReadAll(file)

	var mime string = http.DetectContentType(file_bytes)
	if _, ok = allowed_mime[mime]; !ok {
		return
	}

	var file_id string
	if file_id, err = upload(&file_tee); err != nil {
		return
	}

	created = library_types.NewContent(
		file_id,
		author,
		mime,
		body.Tags,
		body.Featurable,
		body.NSFW,
	)

	err = database.WriteContent(created.Map())
	ok = true
	return
}

func PostContent(request *http.Request) (code int, r_map map[string]interface{}, err error) {
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

	var created library_types.Content
	var author string = request.Context().Value("requester").(string)
	created, ok, err = makeContent(data, file, author)

	switch {
	case err != nil:
		return
	case !ok:
		code = 400
	default:
		code = 200
		r_map = map[string]interface{}{"content": created}
	}

	go database.IncrementPostCount(author)
	return
}
