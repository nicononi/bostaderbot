package repository

import (
	"bostaderbot/repository/database"
	"fmt"

	"github.com/nicononi/collections"
)

type Visited struct {
	ChatId      int64
	Rooms       int
	ApartmentId int
}

// Get Visited Apartments for the given Chat and the amount of rooms
func VisitedForChatAndRooms(chatId int64, rooms int) (collections.List[Visited], error) {
	var visited []Visited
	result := new(collections.SliceList[Visited])

	db := database.GetDB()
	rows, err := db.Query("SELECT chat_id, rooms, apartment_id FROM visited WHERE chat_id=$1 AND rooms=$2", chatId, rooms)

	if err != nil {
		return result, fmt.Errorf("DB.Query: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var vis Visited
		if err := rows.Scan(&vis.ChatId, &vis.Rooms, &vis.ApartmentId); err != nil {
			return nil, fmt.Errorf("fill %q: %v", chatId, err)
		}
		visited = append(visited, vis)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fill %q: %v", chatId, err)
	}

	for _, v := range visited {
		result.Append(v)
	}

	return result, nil
}

// Saves visited apartment
func SaveVisited(chatId int64, rooms int, apartmentId int) error {
	db := database.GetDB()

	insertVisited := "INSERT INTO visited(chat_id, rooms, apartment_id) VALUES($1, $2, $3)"
	_, err := db.Exec(insertVisited, chatId, rooms, apartmentId)

	if err != nil {
		return fmt.Errorf("DB.Exec: %v", err)
	}

	return nil
}
