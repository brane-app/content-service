package main

import (
	"github.com/imonke/monkebase"
	"github.com/imonke/monketype"

	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"testing"
)

const (
	nick  = "user-create"
	email = "user-create@imonke.io"
)

var (
	user  monketype.User
	token string
)

func ferroDelete(url string) (err error) {
	var sendable *http.Request
	if sendable, err = http.NewRequest("DELETE", url, nil); err != nil {
		return
	}

	_, err = ferroRequest(sendable)
	return

}

func newOk(test *testing.T, content monketype.Content) {
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
		test.Errorf("Too many repubs! %d", content.ViewCount)
	}

	if content.CommentCount != 0 {
		test.Errorf("Too many repubs! %d", content.CommentCount)
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

	filePart.Write([]byte("haha yes I am a file"))
	if err = writer.Close(); err != nil {
		panic(err)
	}

	if request, err = http.NewRequest(method, path, bytes.NewReader(body.Bytes())); err != nil {
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

func TestMain(main *testing.M) {
	monkebase.Connect(os.Getenv("MONKEBASE_CONNECTION"))
	user = monketype.NewUser(nick, "", email)

	var err error
	if err = monkebase.WriteUser(user.Map()); err != nil {
		panic(err)
	}

	if token, _, err = monkebase.CreateToken(user.ID); err != nil {
		panic(err)
	}

	var result int = main.Run()
	monkebase.DeleteUser(user.ID)
	os.Exit(result)
}

func Test_postContent(test *testing.T) {
	var set []byte
	var sets [][]byte = [][]byte{
		mustMarshal(map[string]interface{}{
			"mime":       "png",
			"featurable": false,
			"nsfw":       false,
			"tags":       []string{"some", "tags"},
		}),
		mustMarshal(map[string]interface{}{
			"mime":       "png",
			"featurable": false,
			"nsfw":       true,
			"tags":       []string{"some", "tags"},
		}),
		mustMarshal(map[string]interface{}{
			"mime":       "png",
			"featurable": true,
			"nsfw":       false,
			"tags":       []string{"some", "tags"},
		}),
		mustMarshal(map[string]interface{}{
			"mime":       "png",
			"featurable": true,
			"nsfw":       true,
			"tags":       []string{"some", "tags"},
		}),
		mustMarshal(map[string]interface{}{
			"mime":       "png",
			"featurable": true,
			"nsfw":       true,
			"tags":       []string{},
		}),
		mustMarshal(map[string]interface{}{
			"mime":       "png",
			"featurable": true,
			"nsfw":       true,
			"tags":       make([]string, 0),
		}),
		mustMarshal(map[string]interface{}{
			"mime":       "png",
			"featurable": true,
			"nsfw":       true,
			"tags":       nil,
		}),
		mustMarshal(map[string]interface{}{
			"mime": "png",
		}),
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
		request.ParseMultipartForm(MULTIPART_MEM_MAX)

		if code, r_map, err = postContent(request); err != nil {
			test.Fatal(err)
		}

		if code != 200 {
			test.Errorf("got code %d", code)
		}

		var content monketype.Content
		var ok bool
		if content, ok = r_map["content"].(monketype.Content); !ok {
			test.Errorf("%#v", r_map)
		}

		defer ferroDelete(content.FileURL)

		newOk(test, content)
		if content.Author != user.ID {
			test.Errorf("content author mismatch! have: %s, want: %s", content.Author, user.ID)
		}
	}
}

func Test_postContent_badrequest(test *testing.T) {
	var set []byte
	var sets [][]byte = [][]byte{
		mustMarshal(map[string]interface{}{
			"featurable": true,
			"nsfw":       true,
			"tags":       []string{"tag"},
		}),
		mustMarshal(map[string]interface{}{
			"featurable": true,
			"nsfw":       true,
		}),
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
		request.ParseMultipartForm(MULTIPART_MEM_MAX)

		// groudon defaults handles 400 bodies, no testing here
		if code, _, err = postContent(request); err != nil {
			test.Fatal(err)
		}

		if code != 400 {
			test.Errorf("got code %d", code)
		}
	}
}

func Test_rejectBadAuth_blank(test *testing.T) {
	var request *http.Request = new(http.Request)

	var ok bool
	var code int
	var err error
	if _, ok, code, _, err = rejectBadAuth(request); err != nil {
		test.Fatal(err)
	}

	if ok {
		test.Errorf("blank request got through")
	}

	if code != 401 {
		test.Errorf("got code %d", code)
	}
}

func Test_rejectBadAuth_invalid(test *testing.T) {
	var request *http.Request = new(http.Request)
	request.Header = make(http.Header)
	request.Header.Add("Authorization", "foobar")

	var ok bool
	var code int
	var err error
	if _, ok, code, _, err = rejectBadAuth(request); err != nil {
		test.Fatal(err)
	}

	if ok {
		test.Errorf("invalid auth request got through")
	}

	if code != 401 {
		test.Errorf("got code %d", code)
	}
}

func Test_rejectBadAuth_badauth(test *testing.T) {
	var request *http.Request = new(http.Request)
	request.Header = make(http.Header)
	request.Header.Add("Authorization", "Bearer foobar")

	var ok bool
	var code int
	var err error
	if _, ok, code, _, err = rejectBadAuth(request); err != nil {
		test.Fatal(err)
	}

	if ok {
		test.Errorf("bad auth request got through")
	}

	if code != 401 {
		test.Errorf("got code %d", code)
	}
}

func Test_rejectBadAuth(test *testing.T) {
	var request *http.Request = new(http.Request)
	request.Header = make(http.Header)
	request.Header.Add("Authorization", "Bearer "+token)

	var modified *http.Request
	var ok bool
	var err error
	if modified, ok, _, _, err = rejectBadAuth(request); err != nil {
		test.Fatal(err)
	}

	if !ok {
		test.Errorf("ok request did not get through")
	}

	var owner string = modified.Context().Value("requester").(string)
	if owner != user.ID {
		test.Errorf("modified is not owned by %s, but by %s", user.ID, owner)
	}
}

func Test_parseMultipart(test *testing.T) {
	var request *http.Request = mustLocalMultipart("POST", "/", []byte("foobar"))

	var ok bool
	var err error
	if _, ok, _, _, err = parseMultipart(request); err != nil {
		test.Fatal(err)
	}

	if !ok {
		test.Errorf("parseMultipart not ok")
	}

	if request.MultipartForm == nil {
		test.Errorf("request was not multipart parsed")
	}
}

func Test_parseMultipart_err(test *testing.T) {
	var request *http.Request = new(http.Request)
	request.Body = ioutil.NopCloser(strings.NewReader("nice"))

	var ok bool
	var err error
	if _, ok, _, _, err = parseMultipart(request); err != nil {
		test.Fatal(err)
	}

	if ok {
		test.Errorf("parseMultipart is ok on empty request")
	}

	if request.MultipartForm != nil {
		test.Errorf("request was somehow multipart parsed")
	}
}
