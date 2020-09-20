package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
)

var (
	err         error
	debug       bool
	running     bool
	mongoClient *mongo.Client
)

func main() {
	ctx, cancel := newCtx(10)
	defer cancel()
	fmt.Print("MongoDB: Connect ")
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoAuth))
	if err != nil {
		fmt.Println()
		log.Fatal(err)
	}
	fmt.Println("OK")
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	fmt.Print("MongoDB: Ping ")
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println()
		log.Fatal(err)
	}
	fmt.Println("OK")

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
		if running {
			sendMessage(senderID, "Already running!")
			return
		}
		w.WriteHeader(http.StatusOK)
		err := fetchPage()
		if err != nil {
			sendMessage(senderID, err.Error())
		}
		return
	case "/skin":
		sendMessage(senderID, "Please specify a skin code.\n(example: /fetchskin azul-hoodie)")
	case "/test":
		sendPhoto(senderID, "https://rustlabs.com/img/skins/324/39308.png", "New skin released:\n[*Azul Hoodie*](https://rustlabs.com/skin/azul-hoodie)")
	default:
		if strings.HasPrefix(msgText, "/skin ") {
			path := strings.Fields(msgText)
			skin, err := fetchSkin(path[1])
			if err != nil {
				sendMessage(senderID, "Skin not found!")
				fmt.Println(err)
			}
			sendPhoto(senderID, skin.ImageURL, fmt.Sprintf("Skin name: *%v*\nSkin ID: [*%v*](%v)", skin.DisplayName, skin.WorkshopID, skin.PageURL))
			return
		}
		sendMessage(senderID, "Unknown command!")
	}
	w.WriteHeader(http.StatusOK)
}

func fetchPage() error {
	running = true
	resp, err := http.Get("https://rustlabs.com/skins")
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var count int = 1

	var skins []Skin
	fmt.Println("Fetching all skins...")
	doc.Find(".skin-block-2").EachWithBreak(func(i int, s *goquery.Selection) bool {

		// skinName, _ := s.Attr("data-name")
		// new, _ := s.Attr("data-new")
		// isNew := false
		// if new == "NEW" {
		// 	isNew = true
		// }
		// if isNew {
		// 	newSkins = append(newSkins, Skin{DisplayName: skinName, PagePath: pagePath})
		// }
		href, _ := s.Attr("href")
		skin, err := fetchSkin("https:" + href)
		if err != nil {
			fmt.Println(err)
			return true
		}
		if skin.DisplayName == "" || skin.WorkshopID == "" || skin.ItemName == "" {
			fmt.Println("Skipping empty/invalid skin.")
			return true
		}
		// skin.Num = count
		skins = append(skins, skin)
		fmt.Printf("Fetched skin %v/2261 '%v' (%v)\n", count, skin.DisplayName, skin.WorkshopID)
		count++
		// if count == 10 {
		// 	return false
		// }
		return true
	})

	fmt.Println("Skins fetched. Reverse-Upserting them on the database...")

	// var count int64 = 1

	for i := range skins {
		// fmt.Println(skins[len(skins)-1-i])
		skin := skins[len(skins)-1-i]
		skin.Num = i + 1
		err = upsertSkin(skin)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Upserted skin %v/2261 '%v' (%v).\n", skin.Num, skin.DisplayName, skin.WorkshopID)
	}

	fmt.Println("Done!")

	// fmt.Println("New skins:\n", newSkins)
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	running = false
	return nil
}

func fetchSkin(URL string) (Skin, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return Skin{}, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	workshopID := doc.Find(".stats-table").Find("a").Text()
	displayName := doc.Find(".text-column").Children().First().Text()
	imageURL, _ := doc.Find(".icon-column").Find("img").Attr("src")
	imageURL = "https:" + imageURL
	itemPagePath, _ := doc.Find(".tab-block").Find("div").First().Find("a").Attr("href")
	itemPagePath = strings.ReplaceAll(itemPagePath, "/item/", "")
	item, err := getItemByPagePath(itemPagePath)
	if err != nil {
		fmt.Printf("%v (tried fetching item '%v')\n", err, itemPagePath)
		return Skin{}, nil
	}
	return Skin{WorkshopID: workshopID, DisplayName: displayName, PageURL: URL, ImageURL: imageURL, ItemName: item.ItemName}, nil
}

func getItemByPagePath(pagePath string) (item Item, err error) {
	items := mongoClient.Database("rust-skins").Collection("items")
	ctx, cancel := newCtx(5)
	defer cancel()
	err = items.FindOne(ctx, bson.M{"page_path": pagePath}).Decode(&item)
	return item, err
}

// func isInt(s string) bool {
// 	for _, c := range s {
// 		if !unicode.IsDigit(c) {
// 			return false
// 		}
// 	}
// 	return true
// }
