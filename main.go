package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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

	XKCD_CURRENT_URL = "http://xkcd.com/info.0.json"

	RANDOM_INLINE_NUM = 3
)


// ----------------------------------------------------------------------------------
//  types
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


// ----------------------------------------------------------------------------------
//  global variables
// ----------------------------------------------------------------------------------

var (
	telegramToken  = os.Getenv("TELEGRAM_TOKEN")
)


// ----------------------------------------------------------------------------------
//  functions
// ----------------------------------------------------------------------------------

func RandomizeMe(max int) int {
	return rand.Intn(max - 1) + 1
}

// TODO: merge getCurrent and getXkcd into one func
// ----------------------------------------------------------------------------------
func getCurrent() (current int) {

	response, err := http.Get(XKCD_CURRENT_URL)
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
func getXkcd(num int) (picurl string, alt string) {

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
	picurl = obj.Img
	alt = obj.Alt
	return
}


// ----------------------------------------------------------------------------------
//  application entry
// ----------------------------------------------------------------------------------

func main() {
	// print the version number
	log.Println("xkcdbot", VERSION, "#"+GIT_COMMIT, "started")

	// setup seed for random
	rand.Seed(time.Now().Unix())

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
				var results []interface{}

				// add the configured amount of results to the answer
				for i := 0; i < RANDOM_INLINE_NUM; i++ {
					Num := RandomizeMe(getCurrent())
					pUrl, pAlt := getXkcd(Num)
					pic := tgbotapi.NewInlineQueryResultPhotoWithThumb(update.InlineQuery.ID + strconv.Itoa(i), pUrl, pUrl)
					pic.Caption = pAlt
					results = append(results, pic)	
				}
				
				// build the answer
				answer := tgbotapi.InlineConfig{
					InlineQueryID: update.InlineQuery.ID,
					Results:       results,
					CacheTime:     0,
				}
				_, err := bot.AnswerInlineQuery(answer)
				if err != nil {
					log.Println("failed to answer inline query:", err.Error())
				}
			} else {
				// Get latest xkcd
				Num := getCurrent()
				pUrl, pAlt := getXkcd(Num)
				pic := tgbotapi.NewInlineQueryResultPhotoWithThumb(update.InlineQuery.ID, pUrl, pUrl)
				pic.Caption = pAlt

				answer := tgbotapi.InlineConfig{
					InlineQueryID: pic.ID,
					Results:       []interface{}{pic},
					CacheTime:     3,
				}
				_, err := bot.AnswerInlineQuery(answer)
				if err != nil {
					log.Println("failed to answer inline query:", err.Error())
				}
			}
		}
	}
}
