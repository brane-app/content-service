package handlers

import (
	"github.com/brane-app/types-library"

	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"testing"
)

const (
	nick  = "user-create"
	email = "user-create@imonke.io"
)

var (
	pngBytes []byte = []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // png file header
		104, 97, 104, 97, 32, 121, 101, 115, 32, 73, 32, 97, 109, 32, 97, 32, 112, 110, 103, // haha yes I am a png
	}
)

func ferroDelete(url string) (err error) {
	var sendable *http.Request
	if sendable, err = http.NewRequest("DELETE", url, nil); err != nil {
		return
	}

	_, err = ferroRequest(sendable)
	return

}

func newOk(test *testing.T, content types.Content) {
	if content.LikeCount != 0 {
		test.Errorf("Too many likes! %d", content.LikeCount)
	}

	if content.DislikeCount != 0 {
		test.Errorf("Too many dislikes! %d", content.DislikeCount)
	}

	if content.RepubCount != 0 {
		test.Errorf("Too many repubs! %d", content.RepubCount)
	}

	if content.ViewCount != 0 {
		test.Errorf("Too many views! %d", content.ViewCount)
	}

	if content.CommentCount != 0 {
		test.Errorf("Too many comments! %d", content.CommentCount)
	}

	if content.Featured {
		test.Errorf("This new content is featured!")
	}

	if content.Removed {
		test.Errorf("This new content is removed!")
	}
}

func mustLocalMultipart(method, path string, data []byte) (request *http.Request) {
	var body *bytes.Buffer = new(bytes.Buffer)
	var writer *multipart.Writer = multipart.NewWriter(body)

	var dataHeader textproto.MIMEHeader = make(textproto.MIMEHeader)
	dataHeader.Set("Content-Type", "application/json")
	dataHeader.Set("Content-Disposition", `form-data; name="json"`)

	var dataPart io.Writer
	var err error
	if dataPart, err = writer.CreatePart(dataHeader); err != nil {
		panic(err)
	}
	dataPart.Write(data)

	var filePart io.Writer
	if filePart, err = writer.CreateFormFile("file", "file"); err != nil {
		panic(err)
	}

	filePart.Write(pngBytes)
	if err = writer.Close(); err != nil {
		panic(err)
	}

	if request, err = http.NewRequest(method, path, body); err != nil {
		panic(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	return
}

func mustMarshal(it map[string]interface{}) (data []byte) {
	var err error
	if data, err = json.Marshal(it); err != nil {
		panic(err)
	}

	return
}

func Test_PostContent(test *testing.T) {
	var set []byte
	var sets [][]byte = [][]byte{
		mustMarshal(map[string]interface{}{
			"featurable": false,
			"nsfw":       false,
			"tags":       []string{"some", "tags"},
		}),
		mustMarshal(map[string]interface{}{
			"featurable": false,
			"nsfw":       true,
			"tags":       []string{"some", "tags"},
		}),
		mustMarshal(map[string]interface{}{
			"featurable": true,
			"nsfw":       false,
			"tags":       []string{"some", "tags"},
		}),
		mustMarshal(map[string]interface{}{
			"featurable": true,
			"nsfw":       true,
			"tags":       []string{"some", "tags"},
		}),
		mustMarshal(map[string]interface{}{
			"featurable": true,
			"nsfw":       true,
			"tags":       []string{},
		}),
		mustMarshal(map[string]interface{}{
			"featurable": true,
			"nsfw":       true,
			"tags":       make([]string, 0),
		}),
		mustMarshal(map[string]interface{}{
			"featurable": true,
			"nsfw":       true,
			"tags":       nil,
		}),
		mustMarshal(map[string]interface{}{}),
	}

	var request *http.Request
	var code int
	var r_map map[string]interface{}
	var err error

	for _, set = range sets {
		request = mustLocalMultipart("POST", "/", set)
		request = request.WithContext(context.WithValue(
			request.Context(),
			"requester",
			user.ID,
		))
		request.ParseMultipartForm(4 << 2)

		if code, r_map, err = PostContent(request); err != nil {
			test.Fatal(err)
		}

		if code != 200 {
			test.Errorf("got code %d", code)
		}

		var content types.Content
		var ok bool
		if content, ok = r_map["content"].(types.Content); !ok {
			test.Errorf("%#v", r_map)
		}

		defer ferroDelete(content.FileURL)

		newOk(test, content)
		if content.Author != user.ID {
			test.Errorf("content author mismatch! have: %s, want: %s", content.Author, user.ID)
		}
	}
}

func Test_PostContent_badrequest(test *testing.T) {
	var set []byte
	var sets [][]byte = [][]byte{
		[]byte("Why do they call him Donkey Kong if he's a gorilla"),
		make([]byte, 0),
	}

	var request *http.Request
	var code int
	var err error

	for _, set = range sets {
		request = mustLocalMultipart("POST", "/", set)
		request = request.WithContext(context.WithValue(
			request.Context(),
			"requester",
			user.ID,
		))
		request.ParseMultipartForm(4 << 20)

		if code, _, err = PostContent(request); err != nil {
			test.Fatal(err)
		}

		if code != 400 {
			test.Errorf("got code %d", code)
		}
	}
}
