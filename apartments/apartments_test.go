package apartments

import (
	"testing"

	"github.com/nicononi/collections"
)

func TestGetFilteredApartments(t *testing.T) {
	// Create a list of test apartments
	testApartments := new(collections.SliceList[Apartment])
	testApartments.Append(Apartment{Rooms: 3, Normal: true})
	testApartments.Append(Apartment{Rooms: 2, Normal: true})
	testApartments.Append(Apartment{Rooms: 3, Normal: false})

	// Test filtering for apartments with 3 rooms
	filteredApartments := GetFilteredApartments(testApartments, 3)
	if len(filteredApartments.Elements()) != 1 {
		t.Error("Expected to have one element on the filtered apartments list")
	}
	if filteredApartments.Elements()[0].Rooms != 3 || !filteredApartments.Elements()[0].Normal {
		t.Error("Expected the filtered element to have 3 rooms and be normal")
	}

	// Test filtering for apartments with 2 rooms
	filteredApartments = GetFilteredApartments(testApartments, 2)
	if len(filteredApartments.Elements()) != 1 {
		t.Error("Expected to have one element on the filtered apartments list")
	}
	if filteredApartments.Elements()[0].Rooms != 2 || !filteredApartments.Elements()[0].Normal {
		t.Error("Expected the filtered element to have 2 rooms and be normal")
	}

	// Test filtering for apartments with 4 rooms
	filteredApartments = GetFilteredApartments(testApartments, 4)
	if len(filteredApartments.Elements()) != 0 {
		t.Error("Expected to have no elements on the filtered apartments list")
	}
}
