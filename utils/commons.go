package utils

var (
	NUMBER_OF_IMAGES int
	NUMBER_OF_VIDEOS int
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type FactStruct struct {
	Fact string `json:"fact"`
}

type ImageStruct struct {
	URL    string `json:"url"`
	Index  int    `json:"index"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Alt    string `json:"alt"`
}

type VideoStruct struct {
    URL   string `json:"url"`
    Index int    `json:"index"`
    Alt   string `json:"alt"`
}