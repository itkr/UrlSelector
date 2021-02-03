package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/AlecAivazis/survey"
)

type Zoom struct {
	Nanme string `json:"name""`
	URL   string `json:"url""`
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func getResourcePath() (string, error) {
	// TODO: 設定ファイルでpath指定
	if exists("urls.json") {
		return "urls.json", nil
	}
	return "", errors.New("not found")
}

func readJson(resourcePath string) []Zoom {
	bytes, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		log.Fatal(err)
	}
	var zooms []Zoom
	if err := json.Unmarshal(bytes, &zooms); err != nil {
		log.Fatal(err)
	}
	return zooms
}

func openbrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	urls := []string{}
	names := []string{}
	resourcePath, err := getResourcePath()
	if err != nil {
		urls = append(urls, "https://github.com/itkr/UrlSelector/blob/main/README.md")
		names = append(names, "URLの登録方法")
	} else {
		// TODO: goっぽい書き方は？
		zooms := readJson(resourcePath)
		for _, zoom := range zooms {
			urls = append(urls, zoom.URL)
			names = append(names, zoom.Nanme)
		}
	}

	roomIndex := 0
	prompt := &survey.Select{
		Message: "Choose a room:",
		Options: names,
		Default: names[0],
	}
	err = survey.AskOne(prompt, &roomIndex)
	if err != nil {
		log.Fatal(err)
	}

	openbrowser(urls[roomIndex])
}
