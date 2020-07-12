package main

type CreateContentBody struct {
	Mime       string   `json:"mime"`
	NSFW       bool     `json:"nsfw"`
	Featurable bool     `json:"featurable"`
	Tags       []string `json:"tags"`
}

func (_ CreateContentBody) Types() (values map[string]string) {
	values = map[string]string{
		"mime":       "string",
		"nsfw":       "bool",
		"featurable": "bool",
		"tags":       "[]string",
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
