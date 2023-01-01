package main

import (
	"bostaderbot/apartments"
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicononi/collections"
)

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
			msg.Text = handleList(update.Message.CommandArguments())
		case "clean":
			msg.Text = handleClear()
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func handleList(commandArguments string) string {
	rooms, err := strconv.ParseFloat(commandArguments, 32)

	if err != nil {
		return "Could not understand number of rooms"
	}

	allApartments := apartments.GetApartments()
	filteredRooms := apartments.GetFilteredApartments(allApartments, rooms)

	unscannedApartments := new(collections.SliceList[apartments.Apartment])
	scannedIds := new(collections.SliceList[int])

	// Check visited apartments
	for _, v := range filteredRooms.Elements() {
		if !scannedIds.Contains(v.Id) {
			unscannedApartments.Append(v)
		}
	}

	if unscannedApartments.Size() == 0 {
		return "There are no new apartments with " + fmt.Sprintf("%.1f", rooms) + " rooms"
	}

	var result string = "Current new apartments are: \n\n"

	for _, v := range unscannedApartments.Elements() {
		result = result + v.State + ", " + v.Kommun + "\n" + v.Address + "\n" + fmt.Sprintf("%.1f", v.Rooms)
		result = result + " rooms\n" + strconv.Itoa(v.Size) + " sqr\n" + strconv.Itoa(v.Price) + " sek\n" + apartments.BaseUri + v.Uri + "\n\n"
	}

	result = result + "Thanks, Bot"
	return result
}

func handleClear() string {
	return "Work in progress :)"
}

func handleListAll(commandArguments string) string {
	rooms, err := strconv.ParseFloat(commandArguments, 32)

	if err != nil {
		return "Could not understand number of rooms"
	}

	allApartments := apartments.GetApartments()
	filteredRooms := apartments.GetFilteredApartments(allApartments, rooms)

	if filteredRooms.Size() == 0 {
		return "There are no apartments with " + fmt.Sprintf("%.1f", rooms) + " rooms"
	}

	var result string = "Current apartments are: \n\n"

	for _, v := range filteredRooms.Elements() {
		result = result + v.State + ", " + v.Kommun + "\n" + v.Address + "\n" + fmt.Sprintf("%.1f", v.Rooms)
		result = result + " rooms\n" + strconv.Itoa(v.Size) + " sqr\n" + strconv.Itoa(v.Price) + " sek\n" + apartments.BaseUri + v.Uri + "\n\n"
	}

	result = result + "Thanks, Bot"
	return result
}
