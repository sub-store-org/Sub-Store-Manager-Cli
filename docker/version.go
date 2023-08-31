package docker

import (
	"time"

	"sub-store-manager-cli/lib"
	"sub-store-manager-cli/vars"
)

type ReleaseInfo struct {
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

func getBEVersionInfos() (json []ReleaseInfo) {
	_, err := lib.HC.R().SetResult(&json).Get("https://api.github.com/repos/sub-store-org/Sub-Store/releases")
	if err != nil {
		lib.PrintError("Failed to get versions info:", err)
	}
	return
}

func getBEVersionStrs() (v []string) {
	for _, info := range getBEVersionInfos() {
		v = append(v, info.TagName)
	}
	return
}

type FEVersionInfo struct {
	Sha string `json:"sha"`
}

func getFEVersion() (v string) {
	var info FEVersionInfo
	_, err := lib.HC.R().SetResult(&info).Get("https://api.github.com/repos/sub-store-org/Sub-Store-Front-End/commits/master")
	if err != nil {
		lib.PrintError("Failed to get versions info:", err)
	}
	return info.Sha
}

func (c *Container) SetLatestVersion() {
	switch c.ContainerType {
	case vars.ContainerTypeFE:
		c.Version = getFEVersion()[:7]
	case vars.ContainerTypeBE:
		versions := getBEVersionStrs()
		if len(versions) == 0 {
			lib.PrintError("no versions found", nil)
		}
		c.Version = versions[0]
	}
}

func (c *Container) CheckVersionValid() bool {
	if c.ContainerType != vars.ContainerTypeBE {
		return true
	}

	versions := getBEVersionStrs()
	for _, v := range versions {
		if v == c.Version {
			return true
		}
	}
	return false
}
