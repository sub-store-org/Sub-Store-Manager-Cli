package lib

import (
	"errors"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	hcIsInit bool
	HC       *resty.Client
)

func initHttpClient() {
	if hcIsInit {
		return
	}
	client := resty.New()
	HC = client
	hcIsInit = true
}

type LatestRes struct {
	Url         string    `json:"url"`
	AssetsUrl   string    `json:"assets_url"`
	UploadUrl   string    `json:"upload_url"`
	HtmlUrl     string    `json:"html_url"`
	Id          int       `json:"id"`
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	PublishedAt time.Time `json:"published_at"`
	Body        string    `json:"body"`
}

func GetVersionsInfo() []LatestRes {
	var json []LatestRes
	_, err := HC.R().SetResult(&json).Get("https://api.github.com/repos/sub-store-org/Sub-Store/releases")
	if err != nil {
		log.Fatalln("Failed to get versions:", err)
	}
	return json
}

func GetVersionsString() []string {
	var versions []string
	for _, v := range GetVersionsInfo() {
		versions = append(versions, v.TagName)
	}

	return versions
}

func GetLatestVersionString() (string, error) {
	list := GetVersionsString()
	if len(list) == 0 {
		return "", errors.New("no versions found")
	}
	return GetVersionsString()[0], nil
}

func downloadFile(url string, path string) {
	_, err := HC.R().SetOutput(path).Get(url)
	if err != nil {
		log.Fatalln("Failed to download file:", err)
	}
}
