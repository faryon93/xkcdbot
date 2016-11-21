package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
	VERSION    = "0.1.0"
	GIT_COMMIT = "e76c8c01a"

	TELEGRAM_UPDATE_TIMEOUT = 60
)

// ----------------------------------------------------------------------------------
//  global variables
// ----------------------------------------------------------------------------------

type XkcdStruct struct {
	Alt        string `json:"alt"`
	Day        string `json:"day"`
	Img        string `json:"img"`
	Link       string `json:"link"`
	Month      string `json:"month"`
	News       string `json:"news"`
	Num        int    `json:"num"`
	SafeTitle  string `json:"safe_title"`
	Title      string `json:"title"`
	Transcript string `json:"transcript"`
	Year       string `json:"year"`
}

var (
	telegramToken  = os.Getenv("TELEGRAM_TOKEN")
	xkcdCurrentUrl = "http://xkcd.com/info.0.json"
)

// ----------------------------------------------------------------------------------
//  Functions
// ----------------------------------------------------------------------------------
func RandomizeMe(current int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(current-1) + 1
}

// ----------------------------------------------------------------------------------
func getCurrent() (current int) {

	response, err := http.Get(xkcdCurrentUrl)
	if err != nil {
		log.Panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panic(err)
	}
	var obj XkcdStruct
	json.Unmarshal(body, &obj)
	return (obj.Num)
}

// ----------------------------------------------------------------------------------
func getXkcd(num int) (picurl string) {

	url := fmt.Sprint("http://xkcd.com/", num, "/info.0.json")

	response, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panic(err)
	}
	var obj XkcdStruct
	json.Unmarshal(body, &obj)
	return (obj.Img)
}

// ----------------------------------------------------------------------------------
//  application entry
// ----------------------------------------------------------------------------------

func main() {
	// print the version number
	log.Println("xkcdbot", VERSION, "#"+GIT_COMMIT, "started")

	// authenticate with the telegram bot api
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("authorized with telegram bot @%s", bot.Self.UserName)

	// get an update channel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = TELEGRAM_UPDATE_TIMEOUT
	updates, err := bot.GetUpdatesChan(u)

	// process all updates received by the bot
	for update := range updates {
		// we received a private message
		if update.Message != nil {
			// log the request
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// reply to the bot
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello "+update.Message.From.UserName+"! How are you?")
			bot.Send(msg)

			// we received a inline query
		} else if update.InlineQuery != nil {
			log.Printf("[%s] inline query: %s", update.InlineQuery.From.UserName, update.InlineQuery.Query)
			if update.InlineQuery.Query == "random" {
				Num := RandomizeMe(getCurrent())
				pic := tgbotapi.NewInlineQueryResultPhoto(update.InlineQuery.ID, getXkcd(Num))
				pic.ThumbURL = getXkcd(Num)

				answer := tgbotapi.InlineConfig{
					InlineQueryID: update.InlineQuery.ID,
					Results:       []interface{}{pic},
					CacheTime:     0,
				}
				_, err := bot.AnswerInlineQuery(answer)
				if err != nil {
					log.Println("failed to answer inline query:", err.Error())
				}
			}
		}
	}
}
