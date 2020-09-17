package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
)

var (
	debug bool
)

func main() {
	flag.BoolVar(&debug, "d", false, "Debug mode (mostly printing JSON body)")
	flag.Parse()
	http.HandleFunc("/rust-skins-bot", handler)
	if debug {
		fmt.Printf("Debug mode enabled!\n")
	}
	fmt.Printf("Listening on port 1379.\n")
	if err := http.ListenAndServe(":1379", nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	rawJSON, _ := ioutil.ReadAll(r.Body)
	if debug {
		fmt.Println(string(rawJSON))
	}

	// senderName, _ := jsonparser.GetString(rawJSON, "message", "from", "first_name")
	msgText, _ := jsonparser.GetString(rawJSON, "message", "text")
	senderID64, _ := jsonparser.GetInt(rawJSON, "message", "from", "id")
	senderID := strconv.FormatInt(senderID64, 10)
	// msgID64, _ := jsonparser.GetInt(rawJSON, "message", "message_id")
	// msgID := strconv.FormatInt(msgID64, 10)
	fmt.Printf("> %v: %v\n", senderID, msgText)
	switch msgText {
	case "/ping":
		sendMessage(senderID, "Pong!")
	case "/start":
		sendMessage(senderID, "Work In Progress!")
	case "/fetchall":
		err := fetchPage()
		if err != nil {
			sendMessage(senderID, err.Error())
		}
	case "/fetchskin":
		err := fetchSkin("capitan's-ar")
		if err != nil {
			sendMessage(senderID, err.Error())
		}
	case "/test":
		sendPhoto(senderID, "https://rustlabs.com/img/skins/324/39304.png", "New skin released:\n[*No Mercy Kilt*](https://rustlabs.com/skin/no-mercy-kilt)")
	default:
		if strings.HasPrefix(msgText, "/fetchskin ") {
			path := strings.Fields(msgText)
			fetchSkin(path[1])
			return
		}
		sendMessage(senderID, "Unknown command!")
	}
	w.WriteHeader(http.StatusOK)
}

func fetchPage() error {
	resp, err := http.Get("https://rustlabs.com/skins")
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	count := 0

	fmt.Println("fetching latest 5 skins...")

	doc.Find(".skin-block-2").EachWithBreak(func(i int, s *goquery.Selection) bool {

		// For each item found, get the band and title
		href, ok := s.Attr("href")
		if ok {
			fmt.Println(href)
		}
		// band := s.Find("a").Text()
		// title := s.Find("i").Text()
		// fmt.Printf("Review %d: %s - %s\n", i, band, title)
		count++
		if count == 5 {
			return false
		}
		return true
	})

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	return nil
}
func fetchSkin(path string) error {
	resp, err := http.Get("https://rustlabs.com/skin/" + path)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	skinID := doc.Find(".stats-table").Find("tr").Last().Children().Find("a").Text()

	fmt.Println("skin ID: " + skinID)
	// doc.Find("tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
	// 	if s.Children().First().Text() == "Workshop ID" {
	// 		fmt.Println("ID found: " + s.Children().Last().Text())
	// 		return false
	// 	}
	// 	return true
	// })
	return nil
}
