package main

import (
    "log"
    "os"

    "gopkg.in/telegram-bot-api.v4"
)

// ----------------------------------------------------------------------------------
//  constants
// ----------------------------------------------------------------------------------

const (
    VERSION = "0.1.0"
    GIT_COMMIT = "e76c8c01a"

    TELEGRAM_UPDATE_TIMEOUT = 60
)


// ----------------------------------------------------------------------------------
//  global variables
// ----------------------------------------------------------------------------------

var (
    telegramToken = os.Getenv("TELEGRAM_TOKEN")
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
            msg := tgbotapi.NewMessage(update.Message.Chat.ID,"Hello " + update.Message.From.UserName + "! How are you?")
            bot.Send(msg)

        // we received a inline query
        } else if update.InlineQuery != nil {
            log.Printf("[%s] inline query: %s", update.InlineQuery.From.UserName, update.InlineQuery.Query)

            // build the answer to the inline query
            article := tgbotapi.NewInlineQueryResultArticleMarkdown(update.InlineQuery.ID, "Test", "Ich bin die Antwort!")
            article.Description = "Wer bin ich?"
            answer := tgbotapi.InlineConfig{
                InlineQueryID: update.InlineQuery.ID,
                Results: []interface{}{article},
                CacheTime: 0,
            }

            // send the answer
            _, err := bot.AnswerInlineQuery(answer)
            if err != nil {
                log.Println("failed to answer inline query:", err.Error())
            }
        }       
    }
}