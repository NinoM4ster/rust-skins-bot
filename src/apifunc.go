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
