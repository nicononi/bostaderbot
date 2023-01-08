package main

import (
	"bostaderbot/repository/database"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	go startServer()

	var apiKey string

	// If API KEY not present in ENV, check arguments
	fmt.Println("Checking Env API KEY")
	if _, present := os.LookupEnv("API_KEY"); !present {
		fmt.Println("Checking Argument")
		if len(os.Args) == 1 {
			fmt.Println("Need API Key!")
			return
		}

		fmt.Println("Using API KEY from Argument")
		apiKey = os.Args[1]
	} else {
		fmt.Println("Using Env API KEY")
		apiKey = os.Getenv("API_KEY")
	}

	bot, err := tgbotapi.NewBotAPI(apiKey)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /list /latest and /clean."
		case "list":
			msg.Text = handleListAll(update.Message.CommandArguments())
		case "latest":
			msg.Text = handleList(update.Message.Chat.ID, update.Message.CommandArguments())
		case "clean":
			msg.Text = handleClear(update.Message.Chat.ID)
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}

	defer database.GetDB().Close()
}

// This is needed by Cloud Run to keep the Pods running
func startServer() {
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
