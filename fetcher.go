package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const localPath string = "D:\\Practice\\Go\\danbooru\\"
const endpoint string = "http://danbooru.donmai.us/posts.json?limit=1&tags=architecture+no_humans&random=true"
const interval = 10 * time.Minute //10 minutes

var fileName = "random.jpg"

func fetchNewImageLink() string {
	for {
		resp, err := http.Get(endpoint)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)
		var responseBody []interface{}
		if err = dec.Decode(&responseBody); err != nil {
			log.Fatal(err)
		}
		img := responseBody[0].(map[string]interface{})
		if (img["image_height"].(float64) >= 768.0) && (img["image_width"].(float64)/img["image_height"].(float64) >= 1.6) {
			return "http://danbooru.donmai.us" + img["file_url"].(string)
		}
	}
}

func fetchNewImage(fileURL, filePath string) {
	var fo *os.File
	var err error
	if fo, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = fo.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	resp, e := http.Get(fileURL)
	if e != nil {
		log.Fatal(e)
	}
	defer resp.Body.Close()
	_, err = io.Copy(fo, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func wallpaper() {
	fileURL := ""
	for {
		if _, err := os.Stat(localPath + fileName); err == nil {
			err := os.Remove(localPath + fileName)
			if err != nil {
				log.Fatal(err)
			}
		}
		fileURL = fetchNewImageLink()
		fileName = fileURL[strings.LastIndex(fileURL, "/")+1:]
		fetchNewImage(fileURL, localPath+fileName)
		time.Sleep(interval)
	}
}

func userInterface() {
	for {
		fmt.Print("Press q or quit and enter when you want to quit: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if (input[:len(input)-2] == "q") || (input[:len(input)-2] == "quit") {
			break
		}
	}
	if _, err := os.Stat(localPath + fileName); err == nil {
		err := os.Remove(localPath + fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	go wallpaper()
	userInterface()
}
