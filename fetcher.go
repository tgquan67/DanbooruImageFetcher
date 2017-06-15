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
const interval = 600 //in seconds, 600 secconds == 10 minutes

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

func deleteCurrentImg() {
	if _, err := os.Stat(localPath + fileName); err == nil {
		err := os.Remove(localPath + fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func wallpaper(ch chan bool) {
ChangeWallpaper:
	for {
		deleteCurrentImg()
		fileURL := fetchNewImageLink()
		fileName = fileURL[strings.LastIndex(fileURL, "/")+1:]
		fetchNewImage(fileURL, localPath+fileName)
		for i := 0; i < interval; i++ {
			select {
			case <-ch:
				continue ChangeWallpaper
			default:
				time.Sleep(time.Second)
			}
		}
	}
}

func userInterface(ch chan bool) {
WaitForUserInteraction:
	for {
		fmt.Println("Press q or quit and enter when you want to quit.")
		fmt.Println("Press n or next to force next wallpaper right now.")
		fmt.Print("What is your choice? ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		switch input[:len(input)-2] {
		case "q":
			break WaitForUserInteraction
		case "quit":
			break WaitForUserInteraction
		case "n":
			ch <- true
		case "next":
			ch <- true
		}
	}
	deleteCurrentImg()
}

func main() {
	var interrupt = make(chan bool)
	go wallpaper(interrupt)
	userInterface(interrupt)
}
