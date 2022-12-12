package apartments

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"

	"github.com/nicononi/collections"
)

const BaseUri = "https://bostad.stockholm.se"

type Apartment struct {
	Id         int     `json:"AnnonsId"`
	Rooms      float64 `json:"AntalRum"`
	State      string  `json:"Stadsdel"`
	Address    string  `json:"Gatuadress"`
	Kommun     string  `json:"Kommun"`
	Expiration string  `json:"AnnonseradTill"`
	Size       int     `json:"Yta"`
	Price      int     `json:"Hyra"`
	Normal     bool    `json:"Vanlig"`
	Uri        string  `json:"Url"`
}

func GetApartments() collections.List[Apartment] {
	url := BaseUri + "/AllaAnnonser/"

	res, getErr := http.Get(url)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var apartments []Apartment
	jsonErr := json.Unmarshal(body, &apartments)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	result := new(collections.SliceList[Apartment])

	for _, v := range apartments {
		result.Append(v)
	}

	return result
}

func GetFilteredApartments(apartments collections.List[Apartment], roomsNumber float64) collections.List[Apartment] {
	r := new(collections.SliceList[Apartment])
	for _, v := range apartments.Elements() {
		if math.Trunc(v.Rooms) == math.Trunc(roomsNumber) && v.Normal {
			r.Append(v)
		}
	}
	return r
}
