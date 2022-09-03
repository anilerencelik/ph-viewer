package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var WEBHOOK_URL string = "https://webhook.site/ff70e548-7a3d-403b-83e5-d0acea619ede"
var URL string = "https://www.pornhub.com/users/lolloldeneme"
var LASTVIDEOCOUNT = -99
var LASTACCESSLOGIN = "124 years ago"

func sendRequest(url string) *goquery.Document {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer resp.Body.Close()
	httpBody := resp.Body
	node, _ := html.Parse(httpBody)
	document := goquery.NewDocumentFromNode(node)
	return document
}

func parseVideo(document *goquery.Document) int {
	videoNumber := document.Find("ul.subViewsInfoContainer").Find("li").Find("a").Find("span.number").Last().Text()
	intVideoNumber, err := strconv.Atoi(videoNumber)
	if err != nil {
		fmt.Println(err.Error())
		return -1
	}
	return intVideoNumber
}

func parseLogin(document *goquery.Document) string {
	lastLogin := document.Find("dl.moreInformation").Find("dd").First().Next().Next().Text()
	return (lastLogin)
}

func checkNotify(videoNumber int, login string) {
	sendNotify := false
	forVideoCount := false
	LASTACCESSLOGIN = login
	if videoNumber > LASTVIDEOCOUNT {
		LASTVIDEOCOUNT = videoNumber
		sendNotify = true
		forVideoCount = true
	}
	if strings.Contains(login, "seconds ago") {
		sendNotify = true
	}
	if strings.Contains(login, "1 minutes ago") {
		sendNotify = true
	}
	if sendNotify {
		sendNotification(forVideoCount)
	}
}

func sendNotification(isVideoCount bool) {

	yazdir := fmt.Sprintf(`{"videoCount":"%d", "LastLogin":"%s", "isNewVideo": "%t"}`, LASTVIDEOCOUNT, LASTACCESSLOGIN, isVideoCount)
	jsonStr := []byte(yazdir)
	req, err := http.NewRequest("POST", WEBHOOK_URL, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	if !strings.Contains(resp.Status, "200 OK") {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
}

func main() {
	doc := sendRequest(URL)

	uptimeTicker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-uptimeTicker.C:
			checkNotify(parseVideo(doc), parseLogin(doc))
		}
	}
}
