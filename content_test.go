package main

import (
	"github.com/google/uuid"
	"git.gastrodon.io/imonke/monkebase"
	"git.gastrodon.io/imonke/monketype"

	"net/http"
	"os"
	"testing"
)

var (
	content monketype.Content
	id      string = uuid.New().String()
)

func contentOK(test *testing.T, content, target monketype.Content) {
	if content.ID != target.ID {
		test.Errorf("content id mismatch! have: %s, want: %s", content.ID, target.ID)
	}

	if len(content.Tags) != len(target.Tags) {
		test.Errorf("tags have: %v, want: %v", content.Tags, target.Tags)
	}

	var index int
	var tag string
	for index, tag = range target.Tags {
		if tag != content.Tags[index] {
			test.Errorf("tag mismatch at %d! have: %s, want: %s", index, content.Tags[index], tag)
		}
	}
}

func TestMain(main *testing.M) {
	content = monketype.NewContent("", id, "png", nil, true, true)
	monkebase.Connect(os.Getenv("DATABASE_CONNECTION"))

	var result int = main.Run()
	monkebase.DeleteContent(content.ID)
	os.Exit(result)
}

func Test_getContent(test *testing.T) {
	var author, file_url, mime string = uuid.New().String(), "foobar", "png"
	var tags []string = []string{"foo", "bar"}

	var content monketype.Content = monketype.NewContent(file_url, author, mime, tags, false, false)
	var err error
	if err = monkebase.WriteContent(content.Map()); err != nil {
		test.Fatal(err)
	}

	var request *http.Request
	if request, err = http.NewRequest("GET", "/"+content.ID, nil); err != nil {
		test.Fatal(err)
	}

	var code int
	var r_map map[string]interface{}
	if code, r_map, err = getContent(request); err != nil {
		test.Fatal(err)
	}

	if code != 200 {
		test.Errorf("got code %d", code)
	}

	var fetched monketype.Content = monketype.Content{}
	if fetched, err = monketype.ContentFromMap(r_map["content"].(map[string]interface{})); err != nil {
		test.Fatal(err)
	}

	_ = fetched
}

func Test_getContent_notfound(test *testing.T) {
	var request *http.Request
	var err error
	if request, err = http.NewRequest("GET", "https://imonke.io/"+uuid.New().String(), nil); err != nil {
		test.Fatal(err)
	}

	var code int
	var r_map map[string]interface{}
	if code, r_map, err = getContent(request); err != nil {
		test.Fatal(err)
	}

	if code != 404 {
		test.Errorf("got code %d", code)
	}

	if r_map["error"].(string) != "no_such_content" {
		test.Errorf("got response %#v", r_map)
	}
}
