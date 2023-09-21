// Package geocoder provides an easy way to use the Google Geocoding API
package geocoder

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/kelvins/geocoder/structs"
)

// ApiKey The user should set the API KEY provided by Google
var ApiKey string

// ApiUrl The user should set the API URL with the default value being provided by Google.
var ApiUrl = "https://maps.googleapis.com/maps/api/geocode/json?"

// Address structure used in the Geocoding and GeocodingReverse functions
// Note: The FormattedAddress field should be used only for the GeocodingReverse
// to get the formatted address from the Google Geocoding API. It is not used in
// the Geocoding function.
type Address struct {
	Street           string
	Number           int
	Neighborhood     string
	District         string
	City             string
	County           string
	State            string
	Country          string
	PostalCode       string
	FormattedAddress string
	Types            string
}

// Location structure used in the Geocoding and GeocodingReverse functions
type Location struct {
	Latitude  float64
	Longitude float64
}

// Format an address based on the Address structure
// Return the formated address (string)
func (address *Address) FormatAddress() string {

	// Creats a slice with all content from the Address struct
	var content []string
	if address.Number > 0 {
		content = append(content, strconv.Itoa(address.Number))
	}
	content = append(content, address.Street)
	content = append(content, address.Neighborhood)
	content = append(content, address.District)
	content = append(content, address.PostalCode)
	content = append(content, address.City)
	content = append(content, address.County)
	content = append(content, address.State)
	content = append(content, address.Country)

	var formattedAddress string

	// For each value in the content slice check if it is valid
	// and add to the formattedAddress string
	for _, value := range content {
		if value != "" {
			if formattedAddress != "" {
				formattedAddress += ", "
			}
			formattedAddress += value
		}
	}

	return formattedAddress
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
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return results, err
	}

	// The "OK" status indicates that no error has occurred, it means
	// the address was analyzed and at least one geographic code was returned
	if strings.ToUpper(results.Status) != "OK" {
		// If the status is not "OK" check what status was returned
		switch strings.ToUpper(results.Status) {
		case "ZERO_RESULTS":
			err = errors.New("No results found.")
			break
		case "OVER_QUERY_LIMIT":
			err = errors.New("You are over your quota.")
			break
		case "REQUEST_DENIED":
			err = errors.New(results.ErrorMessage)
			break
		case "INVALID_REQUEST":
			err = errors.New(results.ErrorMessage)
			break
		case "UNKNOWN_ERROR":
			err = errors.New("Server error. Please, try again.")
			break
		default:
			break
		}
	}

	return results, err
}

// Geocoding function is used to convert an Address structure
// to a Location structure (latitude and longitude)
func Geocoding(address Address) (Location, error) {

	var location Location

	// Convert whitespaces to +
	formattedAddress := address.FormatAddress()
	formattedAddress = strings.Replace(formattedAddress, " ", "+", -1)

	// Create the URL based on the formated address
	url := ApiUrl + "address=" + formattedAddress

	// Use the API Key if it was set
	if ApiKey != "" {
		url += "&key=" + ApiKey
	}

	// Send the HTTP request and get the results
	results, err := httpRequest(url)
	if err != nil {
		log.Println(err)
		return location, err
	}

	// Get the results (latitude and longitude)
	location.Latitude = results.Results[0].Geometry.Location.Lat
	location.Longitude = results.Results[0].Geometry.Location.Lng

	return location, nil
}

// Convert a structs.Results to a slice of Address structures
func convertResultsToAddress(results structs.Results) (addresses []Address) {

	for index := 0; index < len(results.Results); index++ {
		var address Address

		// Put each component from the AddressComponents slice in the correct field in the Address structure
		for _, component := range results.Results[index].AddressComponents {
			// Check all types of each component
			for _, types := range component.Types {
				switch types {
				case "route":
					address.Street = component.LongName
					break
				case "street_number":
					address.Number, _ = strconv.Atoi(component.LongName)
					break
				case "neighborhood":
					address.Neighborhood = component.LongName
					break
				case "sublocality":
					address.District = component.LongName
					break
				case "sublocality_level_1":
					address.District = component.LongName
					break
				case "locality":
					address.City = component.LongName
					break
				case "administrative_area_level_3":
					address.City = component.LongName
					break
				case "administrative_area_level_2":
					address.County = component.LongName
					break
				case "administrative_area_level_1":
					address.State = component.LongName
					break
				case "country":
					address.Country = component.LongName
					break
				case "postal_code":
					address.PostalCode = component.LongName
					break
				default:
					break
				}
			}
		}

		address.FormattedAddress = results.Results[index].FormattedAddress
		address.Types = results.Results[index].Types[0]

		addresses = append(addresses, address)
	}
	return
}

// GeocodingReverse function is used to convert a Location structure
// to an Address structure
func GeocodingReverse(location Location) ([]Address, error) {

	var addresses []Address

	url := getURLGeocodingReverse(location, "")

	// Send the HTTP request and get the results
	results, err := httpRequest(url)
	if err != nil {
		log.Println(err)
		return addresses, err
	}

	// Convert the results to an Address slice called addresses
	addresses = convertResultsToAddress(results)

	return addresses, nil
}

func GeocodingReverseIntl(location Location, language string) ([]Address, error) {

	var addresses []Address

	url := getURLGeocodingReverse(location, language)

	// Send the HTTP request and get the results
	results, err := httpRequest(url)
	if err != nil {
		log.Println(err)
		return addresses, err
	}

	// Convert the results to an Address slice called addresses
	addresses = convertResultsToAddress(results)

	return addresses, nil
}

func getURLGeocodingReverse(location Location, language string) string {
	// Convert the latitude and longitude from double to string
	latitude := strconv.FormatFloat(location.Latitude, 'f', 8, 64)
	longitude := strconv.FormatFloat(location.Longitude, 'f', 8, 64)

	// Create the URL based on latitude and longitude
	url := ApiUrl + "latlng=" + latitude + "," + longitude

	// Use the API key if it was set
	if ApiKey != "" {
		url += "&key=" + ApiKey
	}

	if language != "" {
		url += "&language=" + language
	}

	return url
}
