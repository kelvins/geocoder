
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
    street: "Avenida Paulista",
    number: 456,
    district: "Bela Vista",
    city: "São Paulo",
    state: "São Paulo",
    country: "Brasil",
    postal_code: "01010-101",
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
  if location.latitude != 0.0 {
    t.Error("Expected: ", 0.0,
            "Received:", location.latitude)
  }
  if location.longitude != 0.0 {
    t.Error("Expected: ", 0.0,
            "Received:", location.longitude)
  }

  address = Address{
    street: "Av. Paulista",
    number: 1578,
    district: "Bela Vista",
    city: "Sao Paulo",
    state: "Sao Paulo",
    country: "Brazil",
    postal_code: "01310-200",
  }

  location, err = Geocoding(address)
  if err != nil {
    t.Error(err)
  }
  if location.latitude != -23.5617633 {
    t.Error("Expected: ", -23.5617633,
            "Received:", location.latitude)
  }
  if location.longitude != -46.6560072 {
    t.Error("Expected: ", -46.6560072,
            "Received:", location.longitude)
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
      latitude: -23.5617633,
      longitude: -46.6560072,
    }

    addresses, err = GeocodingReverse(location)
    if err != nil {
      t.Error(err)
    }
    if len(addresses) > 0 {
      fmt.Println(addresses[0])
    }
}
