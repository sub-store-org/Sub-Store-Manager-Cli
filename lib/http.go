package lib

import (
	"github.com/go-resty/resty/v2"
)

var (
	hcIsInit bool
	HC       *resty.Client
)

func InitHttpClient() {
	if hcIsInit {
		return
	}
	client := resty.New()
	HC = client
	hcIsInit = true
}

func DownloadFile(url string, path string) {
	_, err := HC.R().SetOutput(path).Get(url)
	if err != nil {
		PrintError("Failed to download file:", err)
	}
}
