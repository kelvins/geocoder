
package geocoder

import (
  "fmt"
	"testing"
)

func TestFormatAddress(t *testing.T) {
  var address Address
  formatedAddress := FormatAddress(address)

  if formatedAddress != "" {
    t.Error("Expected: empty string - ",
            "Received:", formatedAddress)
  }

  address = Address{
    street: "Cork Street",
    number: 123,
    city: "Cork",
    country: "Ireland",
  }
  formatedAddress = FormatAddress(address)

  if formatedAddress != "123, Cork Street, Cork, Ireland" {
    t.Error("Expected: 123, Cork Street, Cork, Ireland",
            "Received:", formatedAddress)
  }

  address = Address{
    Street: "Avenida Paulista",
    Number: 456,
    District: "Bela Vista",
    City: "São Paulo",
    State: "São Paulo",
    Country: "Brasil",
    PostalCode: "01010-101",
  }
  formatedAddress = FormatAddress(address)

  if formatedAddress != "456, Avenida Paulista, Bela Vista, 01010-101, São Paulo, São Paulo, Brasil" {
    t.Error("Expected: 456, Avenida Paulista, Bela Vista, 01010-101, São Paulo, São Paulo, Brasil",
            "Received:", formatedAddress)
  }

  address = Address{
    number: 789,
  }
  formatedAddress = FormatAddress(address)

  if formatedAddress != "789" {
    t.Error("Expected: 789",
            "Received:", formatedAddress)
  }
}

func TestGeocoding(t *testing.T) {

  var address Address
  location, err := Geocoding(address)
  if err == nil {
    t.Error("Expected an error")
  }
  if location.Latitude != 0.0 {
    t.Error("Expected: ", 0.0,
            "Received:", location.Latitude)
  }
  if location.Longitude != 0.0 {
    t.Error("Expected: ", 0.0,
            "Received:", location.Longitude)
  }

  address = Address{
    Street: "Av. Paulista",
    Nmber: 1578,
    District: "Bela Vista",
    City: "Sao Paulo",
    State: "Sao Paulo",
    Country: "Brazil",
    PostalCode: "01310-200",
  }

  location, err = Geocoding(address)
  if err != nil {
    t.Error(err)
  }
  if location.Latitude != -23.5617633 {
    t.Error("Expected: ", -23.5617633,
            "Received:", location.Latitude)
  }
  if location.Longitude != -46.6560072 {
    t.Error("Expected: ", -46.6560072,
            "Received:", location.Longitude)
  }
}

func TestGeocodingReverse(t *testing.T) {

    var location Location

    addresses, err := GeocodingReverse(location)
    if err == nil {
      t.Error("Expected an error")
    }
    if len(addresses) > 0 {
      t.Error("Expected no addresses")
    }

    location = Location{
      Latitude: -23.5617633,
      Longitude: -46.6560072,
    }

    addresses, err = GeocodingReverse(location)
    if err != nil {
      t.Error(err)
    }
    if len(addresses) > 0 {
      fmt.Println(addresses[0])
    }
}
