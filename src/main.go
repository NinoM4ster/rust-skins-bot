package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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

	switch msgText {
	case "/ping":
		sendMessage(senderID, "Pong!")
	case "/start":
		sendMessage(senderID, "Work In Progress!")
	default:
		sendMessage(senderID, "Unknown command!")
	}
	w.WriteHeader(http.StatusOK)
}
