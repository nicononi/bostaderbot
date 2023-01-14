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

type Registered struct {
	ChatId int64
	Rooms  int
}

// Get Visited Apartments for the given Chat and the amount of rooms
func VisitedForChat(chatId int64) (collections.List[Visited], error) {
	var visited []Visited
	result := new(collections.SliceList[Visited])

	db := database.GetDB()
	rows, err := db.Query("SELECT chat_id, rooms, apartment_id FROM visited WHERE chat_id=$", chatId)

	if err != nil {
		return nil, fmt.Errorf("DB.Query: %v", err)
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

func RegisterForUpdates(chatId int64) error {
	db := database.GetDB()

	insertRegistered := "INSERT INTO registered(chat_id) VALUES($1)"
	_, err := db.Exec(insertRegistered, chatId)

	if err != nil {
		return fmt.Errorf("DB.Exec: %v", err)
	}

	return nil
}

func GetRegisteredChats() (collections.List[Registered], error) {
	var registered []Registered
	result := new(collections.SliceList[Registered])

	db := database.GetDB()
	rows, err := db.Query(`SELECT v.chat_id, v.rooms 
	FROM registered AS r 
	INNER JOIN visited AS v on r.chat_id = v.chat_id
	GROUP BY v.chat_id, v.rooms`)

	if err != nil {
		return nil, fmt.Errorf("DB.Query: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var reg Registered
		if err := rows.Scan(&reg.ChatId, &reg.Rooms); err != nil {
			return nil, fmt.Errorf("fill : %v", err)
		}
		registered = append(registered, reg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("fill %v", err)
	}

	for _, v := range registered {
		result.Append(v)
	}

	return result, nil
}
