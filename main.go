package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	// 0. 引数指定
	flag.Parse()
	fileName := flag.Arg(0)
	if exists(fileName) {
		return fileName, nil
	}
	// 1. 実行ファイルの近くのjson
	fileName = "UrlSelector.json"
	if exists(fileName) {
		return fileName, nil
	}
	// 2. 環境変数
	fileName = os.Getenv("URLSELECTOR_CONFIG")
	if exists(fileName) {
		return fileName, nil
	}
	// 3. ホームディレクトリのjson
	switch runtime.GOOS {
	case "linux":
		fileName = filepath.Join(os.Getenv("HOME"), "UrlSelector.json")
	case "windows":
		fileName = filepath.Join(os.Getenv("HOMEPATH"), "UrlSelector.json")
	case "darwin":
		fileName = filepath.Join(os.Getenv("HOME"), "UrlSelector.json")
	}
	if exists(fileName) {
		return fileName, nil
	}
	// 4. エラー
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
	fullPath, _ := filepath.Abs(resourcePath)
	fmt.Printf("Config file path: %s\n", fullPath)

	roomIndex := 0
	prompt := &survey.Select{
		Message: "Choose a title:",
		Options: names,
		Default: names[0],
	}
	err = survey.AskOne(prompt, &roomIndex)
	if err != nil {
		log.Fatal(err)
	}

	openbrowser(urls[roomIndex])
}
