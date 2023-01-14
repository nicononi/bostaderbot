package commands

import (
	"bostaderbot/apartments"
	"bostaderbot/repository"
	"fmt"
	"strconv"

	"github.com/nicononi/collections"
)

// Function to handle the Delta command
func HandleListDeltaCmd(chatId int64, commandArguments string) string {
	rooms, err := strconv.Atoi(commandArguments)

	if err != nil {
		return "Could not understand number of rooms"
	}

	return HandleListDelta(chatId, rooms)
}

// Function to handle the Delta command
func HandleListDelta(chatId int64, rooms int) string {
	allApartments := apartments.GetApartments()
	filteredRooms := apartments.GetFilteredApartments(allApartments, rooms)

	unscannedApartments := new(collections.SliceList[apartments.Apartment])

	results, err := repository.VisitedForChatAndRooms(chatId, rooms)
	if err != nil {
		fmt.Println(err)
		return "Could not obtain for " + fmt.Sprintf("%d", rooms) + " rooms. Try again later"
	}

	// Keep apartment IDs to make the comparison
	scannedIds := new(collections.SliceList[int])
	for _, v := range results.Elements() {
		scannedIds.Append(v.ApartmentId)
	}

	// Check visited apartments
	for _, v := range filteredRooms.Elements() {
		if !scannedIds.Contains(v.Id) {
			unscannedApartments.Append(v)
		}
	}

	// Check if we have unvisited apartments
	if unscannedApartments.Size() == 0 {
		return "There are no new apartments with " + fmt.Sprintf("%d", rooms) + " rooms"
	}

	var result string = "Current new apartments are: \n\n"

	for _, v := range unscannedApartments.Elements() {
		//Save apartment into the DB
		err := repository.SaveVisited(chatId, rooms, v.Id)

		if err != nil {
			fmt.Println(err)
			return "Could not save new Apartments for " + fmt.Sprintf("%d", rooms) + " rooms. Try again later"
		}

		//Append apartment to result
		result = result + v.State + ", " + v.Kommun + "\n" + v.Address + "\n" + fmt.Sprintf("%.1f", v.Rooms)
		result = result + " rooms\n" + strconv.Itoa(v.Size) + " sqr\n" + strconv.Itoa(v.Price) + " sek\n" + apartments.BaseUri + v.Uri + "\n\n"
	}

	result = result + "Thanks, Bot"
	return result
}
