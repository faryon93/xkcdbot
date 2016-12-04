package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/faryon93/xkcdbot/xkcd"

	"gopkg.in/telegram-bot-api.v4"
)

// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
	VERSION    = "0.1.1"
	GIT_COMMIT = "e76c8c01a"

	TELEGRAM_UPDATE_TIMEOUT = 60
	RANDOM_INLINE_NUM = 3
)


// ----------------------------------------------------------------------------------
//  global variables
// ----------------------------------------------------------------------------------

var (
	telegramToken  = os.Getenv("TELEGRAM_TOKEN")
)


// ----------------------------------------------------------------------------------
//  application entry
// ----------------------------------------------------------------------------------

func main() {
	// print the version number
	log.Println("xkcdbot", VERSION, "#" + GIT_COMMIT, "started")

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
				// fetch the latest comic, in order
				// to get the interval for randomizing
				latest, err := xkcd.GetComic(xkcd.CURRENT_COMIC)
				if err != nil {
					log.Println("failed to fetch latest comic:", err.Error())
					continue
				}

				// add the configured amount of results to the answer
				var results []interface{}
				for i := 0; i < RANDOM_INLINE_NUM; i++ {
					// TODO: comic 404 is not defined, handle that...
					num := RandomizeMe(latest.Num)
					comic, err := xkcd.GetComic(num)
					if err != nil {
						log.Printf("failed to fetch comic %d: %s", num, err.Error())
						continue
					}

					// add the new result
					resultId := update.InlineQuery.ID + strconv.Itoa(i)
					pic := tgbotapi.NewInlineQueryResultPhotoWithThumb(resultId, comic.Img, comic.Img)
					pic.Caption = comic.Alt
					results = append(results, pic)	
				}
				
				// build the answer and send to client
				answer := tgbotapi.InlineConfig{
					InlineQueryID: update.InlineQuery.ID,
					Results:       results,
					CacheTime:     0,
				}
				_, err = bot.AnswerInlineQuery(answer)
				if err != nil {
					log.Println("failed to answer inline query:", err.Error())
				}
			} else {
				// Get latest xkcd
				comic, err := xkcd.GetComic(xkcd.CURRENT_COMIC)
				if err != nil {
					log.Println("failed to fetch current comic:", err.Error())
					continue
				}

				// build the only inline result
				pic := tgbotapi.NewInlineQueryResultPhotoWithThumb(update.InlineQuery.ID, comic.Img, comic.Img)
				pic.Caption = comic.Alt

				// send the answer with results to the client
				answer := tgbotapi.InlineConfig{
					InlineQueryID: pic.ID,
					Results:       []interface{}{pic},
					CacheTime:     3,
				}
				_, err = bot.AnswerInlineQuery(answer)
				if err != nil {
					log.Println("failed to answer inline query:", err.Error())
				}
			}
		}
	}
}


// ----------------------------------------------------------------------------------
//  helper functions
// ----------------------------------------------------------------------------------

func RandomizeMe(max int) int {
	return rand.Intn(max - 1) + 1
}