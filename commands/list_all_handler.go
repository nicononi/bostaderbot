package commands

import (
	"bostaderbot/apartments"
	"fmt"
	"strconv"
)

func HandleListAllCmd(commandArguments string) string {
	rooms, err := strconv.Atoi(commandArguments)

	if err != nil {
		return "Could not understand number of rooms"
	}

	allApartments := apartments.GetApartments()
	filteredRooms := apartments.GetFilteredApartments(allApartments, rooms)

	if filteredRooms.Size() == 0 {
		return "There are no apartments with " + fmt.Sprintf("%d", rooms) + " rooms"
	}

	var result string = "Current apartments are: \n\n"

	for _, v := range filteredRooms.Elements() {
		result = result + v.State + ", " + v.Kommun + "\n" + v.Address + "\n" + fmt.Sprintf("%.1f", v.Rooms)
		result = result + " rooms\n" + strconv.Itoa(v.Size) + " sqr\n" + strconv.Itoa(v.Price) + " sek\n" + apartments.BaseUri + v.Uri + "\n\n"
	}

	result = result + "Thanks, Bot"
	return result
}
