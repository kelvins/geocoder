package geocoder

import (
	"errors"
	"testing"
)

func TestFormatAddress(t *testing.T) {
	var address1 Address
	var address2 Address
	var address3 Address
	var address4 Address

	address2 = Address{
		Street:  "Cork Street",
		Number:  123,
		City:    "Cork",
		Country: "Ireland",
	}

	address3 = Address{
		Street:     "Avenida Paulista",
		Number:     456,
		District:   "Bela Vista",
		City:       "São Paulo",
		State:      "São Paulo",
		Country:    "Brasil",
		PostalCode: "01010-101",
	}

	address4 = Address{
		Number: 789,
	}

	// Table tests
	var tTests = []struct {
		address         Address
		formatedAddress string
	}{
		{address1, ""},
		{address2, "123, Cork Street, Cork, Ireland"},
		{address3, "456, Avenida Paulista, Bela Vista, 01010-101, São Paulo, São Paulo, Brasil"},
		{address4, "789"},
	}

	// Test with all values from the tTests
	for _, pair := range tTests {
		formatedAddress := FormatAddress(pair.address)

		if formatedAddress != pair.formatedAddress {
			t.Error("Expected:", pair.formatedAddress,
				"Received:", formatedAddress)
		}
	}
}

func TestGeocoding(t *testing.T) {

	var address1 Address
	var address2 Address

	var location1 Location
	var location2 Location

	location1 = Location{
		Latitude:  0.0,
		Longitude: 0.0,
	}

	location2 = Location{
		Latitude:  -23.5617633,
		Longitude: -46.6560072,
	}

	address2 = Address{
		Street:     "Av. Paulista",
		Number:     1578,
		District:   "Bela Vista",
		City:       "Sao Paulo",
		State:      "Sao Paulo",
		Country:    "Brazil",
		PostalCode: "01310-200",
	}

	// Table tests
	var tTests = []struct {
		address  Address
		location Location
		err      error
	}{
		{address1, location1, errors.New("Empty Address")},
		{address2, location2, nil},
	}

	// Test with all values from the tTests
	for _, pair := range tTests {
		location, err := Geocoding(pair.address)

		if pair.err != nil {
			if err == nil {
				t.Error("Expected:", pair.err,
					"Received: nil")
			}
		} else {
			if err != nil {
				t.Error("Expected: nil",
					"Received:", err)
			}
		}
		if location.Latitude != pair.location.Latitude {
			t.Error("Expected:", pair.location.Latitude,
				"Received:", location.Latitude)
		}
		if location.Longitude != pair.location.Longitude {
			t.Error("Expected:", pair.location.Longitude,
				"Received:", location.Longitude)
		}
	}
}

func TestGeocodingReverse(t *testing.T) {

	var location1 Location
	var location2 Location

	var address1 Address
	var address2 Address

	location2 = Location{
		Latitude:  -23.5617633,
		Longitude: -46.6560072,
	}

	address2 = Address{
		Street:  "Avenida Paulista",
		County:  "São Paulo",
		State:   "São Paulo",
		Country: "Brazil",
	}

	// Table tests
	var tTests = []struct {
		location Location
		address  Address
		err      error
	}{
		{location1, address1, errors.New("Empty Location")},
		{location2, address2, nil},
	}

	// Test with all values from the tTests
	for _, pair := range tTests {
		addresses, err := GeocodingReverse(pair.location)

		if pair.err != nil {
			if err == nil {
				t.Error("Expected:", pair.err,
					"Received: nil")
			}
		} else {
			if err != nil {
				t.Error("Expected: nil",
					"Received:", err)
			} else if len(addresses) > 0 {
				if addresses[0].Street != pair.address.Street {
					t.Error("Expected Street:", pair.address.Street,
						"Received Street:", addresses[0].Street)
				}
				if addresses[0].County != pair.address.County {
					t.Error("Expected City:", pair.address.County,
						"Received City:", addresses[0].County)
				}
				if addresses[0].State != pair.address.State {
					t.Error("Expected State:", pair.address.State,
						"Received State:", addresses[0].State)
				}
				if addresses[0].Country != pair.address.Country {
					t.Error("Expected Country:", pair.address.Country,
						"Received Country:", addresses[0].Country)
				}
			} else {
				t.Error("Expected at least 1 result Received 0 results:")
			}
		}
	}
}
