package geocoder

import (
	"errors"
	"os"
	"reflect"
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
		address          Address
		formattedAddress string
	}{
		{address1, ""},
		{address2, "123, Cork Street, Cork, Ireland"},
		{address3, "456, Avenida Paulista, Bela Vista, 01010-101, São Paulo, São Paulo, Brasil"},
		{address4, "789"},
	}

	// Test with all values from the tTests
	for _, pair := range tTests {
		formattedAddress := pair.address.FormatAddress()

		if formattedAddress != pair.formattedAddress {
			t.Error("Expected:", pair.formattedAddress,
				"Received:", formattedAddress)
		}
	}
}

func TestGeocoding(t *testing.T) {

	ApiKey = os.Getenv("API_KEY")

	var address1 Address
	var address2 Address

	var location1 Location
	var location2 Location

	location1 = Location{
		Latitude:  0.0,
		Longitude: 0.0,
	}

	location2 = Location{
		Latitude:  -23.5615171,
		Longitude: -46.655961,
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

	ApiKey = os.Getenv("API_KEY")

	var location1 Location
	var location2 Location

	var address1 Address
	var address2 Address

	location2 = Location{
		Latitude:  -23.5617633,
		Longitude: -46.6560072,
	}

	address2 = Address{
		Street:   "Avenida Paulista",
		Number:   1540,
		District: "Bela Vista",
		County:   "São Paulo",
		State:    "São Paulo",
		Country:  "Brazil",
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
				if addresses[0].Number != pair.address.Number {
					t.Error("Expected Number:", pair.address.Number,
						"Received Number:", addresses[0].Number)
				}
				if addresses[0].District != pair.address.District {
					t.Error("Expected District:", pair.address.District,
						"Received District:", addresses[0].District)
				}
				if addresses[0].City != pair.address.City {
					t.Error("Expected City:", pair.address.City,
						"Received City:", addresses[0].City)
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
				t.Error("Expected at least 1 result Received 0 results")
			}
		}
	}
}

func TestHttpRequest(t *testing.T) {

	ApiKey = os.Getenv("API_KEY")

	// Table tests
	var tTests = []struct {
		url string
		err error
	}{
		{"", errors.New("URL invalid")},
		{"https://api.sunrise-sunset.org/json?lat=36.7201600&lng=-4.4203400", errors.New("JSON invalid")},
	}

	// Test with all values from the tTests
	for _, pair := range tTests {
		_, err := httpRequest(pair.url)

		if pair.err != nil && err == nil {
			t.Error("Expected error:", pair.err, "Received: nil")
		}
	}
}

func TestWithInvalidApiKey(t *testing.T) {
	ApiKey = "0123456789abcdefghijklmnopqrstuvxyz"

	address := Address{
		Street:  "Cork Street",
		Number:  123,
		City:    "Cork",
		Country: "Ireland",
	}

	// Table tests
	var tTests1 = []struct {
		address Address
		err     error
	}{
		{address, errors.New("Request Denied")},
	}

	// Test with all values from the tTests
	for _, pair := range tTests1 {
		_, err := Geocoding(pair.address)

		if err == nil {
			t.Error("Expected:", pair.err, "Received: nil")
		}
	}

	location := Location{
		Latitude:  -23.5617633,
		Longitude: -46.6560072,
	}

	// Table tests
	var tTests2 = []struct {
		location Location
		err      error
	}{
		{location, errors.New("Request Denied")},
	}

	// Test with all values from the tTests
	for _, pair := range tTests2 {
		_, err := GeocodingReverse(pair.location)

		if err == nil {
			t.Error("Expected:", pair.err, "Received: nil")
		}
	}
}

func TestGeocodingReverseIntl(t *testing.T) {
	ApiKey = os.Getenv("API_KEY")
	type args struct {
		location Location
		language string
	}
	tests := []struct {
		name    string
		args    args
		want    Address
		wantErr bool
	}{
		{
			name: "Returns address in Ukrainian",
			args: args{
				Location{50.006414, 36.252432},
				"uk",
			},
			want: Address{
				Street:           "Лермонтовська вулиця",
				Number:           7,
				Neighborhood:     "",
				District:         "Київський район",
				City:             "Харківська міськрада",
				County:           "",
				State:            "Харківська область",
				Country:          "Україна",
				PostalCode:       "61000",
				FormattedAddress: "Лермонтовська вулиця, 7, Харків, Харківська область, Україна, 61000",
				Types:            "street_address",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeocodingReverseIntl(tt.args.location, tt.args.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeocodingReverseIntl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got[0], tt.want) {
				t.Errorf("GeocodingReverseIntl() = %v, want %v", got, tt.want)
			}
		})
	}
}
