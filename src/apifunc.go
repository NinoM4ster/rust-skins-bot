package main

import (
	"net/http"
	"net/url"
)

const (
	// token string = "YOUR_BOT_TOKEN_HERE"
	api string = "https://api.telegram.org/bot" + token
)

func sendMessage(userID, message string) error {
	params := url.Values{}
	params.Add("chat_id", userID)
	params.Add("text", message)
	_, err := http.Get(api + "/sendMessage?" + params.Encode())
	return err
}

func sendPhoto(userID, fileURL, caption string) error {
	params := url.Values{}
	params.Add("chat_id", userID)
	params.Add("photo", fileURL)
	params.Add("parse_mode", "MarkdownV2")
	if len(caption) > 0 {
		params.Add("caption", caption)
	}
	_, err := http.Get(api + "/sendPhoto?" + params.Encode())
	return err
}
