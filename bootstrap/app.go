package bootstrap

import (
	"encoding/json"
	"fmt"

	"github.com/alainQtec/FileServer/v3/pkg/conf"
	"github.com/alainQtec/FileServer/v3/pkg/request"
	"github.com/alainQtec/FileServer/v3/pkg/util"
	"github.com/hashicorp/go-version"
)

// InitApplication ie: Initialize application constants
func InitApplication() {
	fmt.Print(`
   ___ _                 _
  / __\ | ___  _   _  __| |_ __ _____   _____
 / /  | |/ _ \| | | |/ _  | '__/ _ \ \ / / _ \
/ /___| | (_) | |_| | (_| | | |  __/\ V /  __/
\____/|_|\___/ \__,_|\__,_|_|  \___| \_/ \___|

   V` + conf.BackendVersion + `  Commit #` + conf.LastCommit + `  Pro=` + conf.IsPro + `
================================================

`)
	go CheckUpdate()
}

type GitHubRelease struct {
	URL  string `json:"html_url"`
	Name string `json:"name"`
	Tag  string `json:"tag_name"`
}

// Check for updates
func CheckUpdate() {
	client := request.NewClient()
	res, err := client.Request("GET", "https://api.github.com/repos/FileServer/FileServer/releases", nil).GetResponse()
	if err != nil {
		util.Log().Warning("Update check failed, %s", err)
		return
	}

	var list []GitHubRelease
	if err := json.Unmarshal([]byte(res), &list); err != nil {
		util.Log().Warning("Update check failed, %s", err)
		return
	}

	if len(list) > 0 {
		present, err1 := version.NewVersion(conf.BackendVersion)
		latest, err2 := version.NewVersion(list[0].Tag)
		if err1 == nil && err2 == nil && latest.GreaterThan(present) {
			util.Log().Info("there is a new version [%s] available , 下载：%s", list[0].Name, list[0].URL)
		}
	}

}
