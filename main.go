package main

import (
	"bostaderbot/commands"
	"bostaderbot/repository"
	"bostaderbot/repository/database"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	bot  *tgbotapi.BotAPI
	once sync.Once
)

// getBot lazily instantiates a TelegamBOT once. Users of Cloud Run or
// Cloud Functions may wish to skip this lazy instantiation and connect as soon
// as the function is loaded. This is primarily to help testing.
func getBot(apiKey string) *tgbotapi.BotAPI {
	once.Do(func() {
		var err error
		bot, err = tgbotapi.NewBotAPI(apiKey)

		if err != nil {
			log.Panic(err)
		}
	})
	return bot
}

func main() {
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

	bot := getBot(apiKey)
	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go startServer()

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
			msg.Text = commands.HandleListAllCmd(update.Message.CommandArguments())
		case "latest":
			msg.Text = commands.HandleListDeltaCmd(update.Message.Chat.ID, update.Message.CommandArguments())
		case "register":
			msg.Text = commands.HandleRegisterCmd(update.Message.Chat.ID)
		case "clean":
			msg.Text = commands.HandleClearCmd(update.Message.Chat.ID)
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
	http.HandleFunc("/reminder", notify)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func notify(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			registeredChats, err := repository.GetRegisteredChats()

			if err != nil {
				fmt.Println(err)
				fmt.Fprintf(w, "Could not notify users, try again later!")
				return
			}

			for _, v := range registeredChats.Elements() {
				// Create a new MessageConfig. We don't have text yet,
				// so we leave it empty.
				msg := tgbotapi.NewMessage(v.ChatId, "")
				msg.Text = commands.HandleListDelta(v.ChatId, v.Rooms)
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}

			fmt.Fprintf(w, "Users notified!")
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
