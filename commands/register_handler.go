package commands

import (
	"bostaderbot/repository"
	"fmt"
)

func HandleRegisterCmd(chatId int64) string {
	err := repository.RegisterForUpdates(chatId)
	if err != nil {
		fmt.Println(err)
		return "Could not register for updates. Please try again."
	}

	return "Registered for updates. You will start receiving a daily message starting tomorrow."
}
