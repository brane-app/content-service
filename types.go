package main

import (
	"github.com/gastrodon/groudon/v2"
)

type CreateContentBody struct {
	Mime       string   `json:"mime"`
	NSFW       bool     `json:"nsfw"`
	Featurable bool     `json:"featurable"`
	Tags       []string `json:"tags"`
}

func (_ CreateContentBody) Validators() (values map[string]func(interface{}) (bool, error)) {
	values = map[string]func(interface{}) (bool, error){
		"nsfw":       groudon.ValidBool,
		"featurable": groudon.ValidBool,
		"tags":       groudon.ValidStringSlice,
	}

	return
}

func (_ CreateContentBody) Defaults() (values map[string]interface{}) {
	values = map[string]interface{}{
		"nsfw":       false,
		"featurable": true,
		"tags":       make([]string, 0),
	}

	return
}
