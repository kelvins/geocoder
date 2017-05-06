// Package geocoder
package geocoder

import (
	"log"
	"errors"
	"strings"
	"strconv"
	"net/http"
	"encoding/json"

	"github.com/kelvins/geocoder/structs"
)

var geocodeApiUrl string
var Key string

// Define the Geocode API URL
func init() {
	geocodeApiUrl = "https://maps.googleapis.com/maps/api/geocode/json?"
}

// Address structure used in the Geocoding and GeocodingReverse functions
type Address struct {
	street      string
	number      int
	district    string
	city        string
	county      string
	state       string
	country     string
	postal_code string
	types       string
}

// Location structure used in the Geocoding and GeocodingReverse functions
type Location struct {
   latitude  float64
   longitude float64
}

// Format an address based on the Address structure
// Return the formated address (string)
func FormatAddress(address Address) string {

	// Creats a slice with all content from the Address struct
	var content []string
	if address.number > 0 {
		content = append(content, strconv.Itoa(address.number))
	}
	content = append(content, address.street)
	content = append(content, address.district)
	content = append(content, address.postal_code)
	content = append(content, address.city)
	content = append(content, address.county)
	content = append(content, address.state)
	content = append(content, address.country)

	var formatedAddress string
	// For each value in the content slice check if it is valid
	// and add to the formatedAddress string
	for index := 0; index < len(content); index++ {
		if content[index] != "" {
			if formatedAddress != "" {
				formatedAddress += ", "
			}
			formatedAddress += content[index]
		}
	}

	return formatedAddress
}

// httpRequest function send the HTTP request, decode the JSON
// and return a Results structure
func httpRequest(url string) (structs.Results, error) {

	var results structs.Results

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return results, err
	}

	// For control over HTTP client headers, redirect policy, and other settings, create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		return results, err
	}

	// Callers should close resp.Body when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Use json.Decode for reading streams of JSON data
	err = json.NewDecoder(resp.Body).Decode(&results);
	if err != nil {
		return results, err
	}

	// Check if the status returned is OK
	if strings.ToLower(results.Status) != "ok" {
		err = errors.New("Status is not OK")
		return results, err
	}

	// Check if we have some result to get
	if len(results.Results) == 0 {
		err = errors.New("No results found")
		return results, err
	}

	return results, nil
}

// Geocoding function is used to convert an Address structure to a
// Location structure (latitude and longitude)
func Geocoding(address Address) (Location, error) {

	var location Location

	// Convert whitespaces to +
	formatedAddress := FormatAddress(address)
	formatedAddress = strings.Replace(formatedAddress, " ", "+", -1)

	// Create the URL based on the formated address
	url := geocodeApiUrl + "address=" + formatedAddress
	if Key != "" {
		url += "&key=" + Key
	}

	results, err := httpRequest(url)
	if err != nil {
		log.Println(err)
		return location, err
	}

	// Get the results (latitude and longitude)
	location.latitude  = results.Results[0].Geometry.Location.Lat;
	location.longitude = results.Results[0].Geometry.Location.Lng;

	return location, nil
}

// Convert a structs.Results to a slice of Address structures
func convertResultsToAddress(results structs.Results) (addresses []Address) {

	for index := 0; index < len(results.Results); index++ {
		var address Address

		for component := 0; component < len(results.Results[index].AddressComponents); component++ {

			switch results.Results[index].AddressComponents[component].Types[0] {
			case "route":
				address.street = results.Results[index].AddressComponents[component].LongName
				break
			case "street_number":
				address.number, _ = strconv.Atoi(results.Results[index].AddressComponents[component].LongName)
				break
			case "locality":
				address.district = results.Results[index].AddressComponents[component].LongName
				break
			case "administrative_area_level_3":
				address.city = results.Results[index].AddressComponents[component].LongName
				break
			case "administrative_area_level_2":
				address.county = results.Results[index].AddressComponents[component].LongName
				break
			case "administrative_area_level_1":
				address.state = results.Results[index].AddressComponents[component].LongName
				break
			case "country":
				address.country = results.Results[index].AddressComponents[component].LongName
				break
			case "postal_code":
				address.postal_code = results.Results[index].AddressComponents[component].LongName
				break
			default:
				break
			}
		}

		address.types = results.Results[index].Types[0]

		addresses = append(addresses, address)
	}
	return
}

// GeocodingReverse function is used to convert a Location structure
// to an Address structure
func GeocodingReverse(location Location) ([]Address, error) {

	var addresses []Address

	// Create the URL based on the formated address
	latitude  := strconv.FormatFloat(location.latitude, 'f', 8, 64)
	longitude := strconv.FormatFloat(location.longitude, 'f', 8, 64)

	url := geocodeApiUrl + "latlng=" + latitude + "," + longitude
	if Key != "" {
		url += "&key=" + Key
	}

	results, err := httpRequest(url)
	if err != nil {
		log.Println(err)
		return addresses, err
	}

	addresses = convertResultsToAddress(results)

	return addresses, nil
}
