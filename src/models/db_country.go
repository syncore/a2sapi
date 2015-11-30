package models

// db_country.go - Model for country information returned by database

// DbCountry This struct is for the JSON representation displayed by the API
type DbCountry struct {
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	Continent   string `json:"region"`
	State       string `json:"state"`
}
