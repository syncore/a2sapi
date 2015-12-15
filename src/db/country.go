package db

// country.go - Country geolocation database lookup.

import (
	"net"
	"steamtest/src/constants"
	"steamtest/src/logger"
	"steamtest/src/models"

	"github.com/oschwald/maxminddb-golang"
)

// This is an intermediate struct to represent the MaxMind DB format, not for JSON
type mmdbformat struct {
	Country struct {
		IsoCode string            `maxminddb:"iso_code"`
		Names   map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	Continent struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"continent"`
	Subdivisions []struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"subdivisions"`
}

func getDefaultCountryData() *models.DbCountry {
	return &models.DbCountry{
		CountryName: "Unknown",
		CountryCode: "Unknown",
		Continent:   "Unknown",
		State:       "Unknown",
	}
}

// OpenCountryDB opens the country lookup database for reading. The caller of
// this function will be responsinble for calling a .Close() on the Reader pointer
// returned by this function.
func OpenCountryDB() (*maxminddb.Reader, error) {
	// Note: the caller of this function needs to handle db.Close()
	db, err := maxminddb.Open(constants.CountryDbFilePath)
	if err != nil {
		logger.LogAppError(err)
		panic("Unable to open country database!")
	}
	return db, nil
}

// GetCountryInfo attempts to retrieve the country information for a given IP,
// returning the result as a country model object over the corresponding result channel.
func GetCountryInfo(ch chan<- *models.DbCountry, db *maxminddb.Reader, ipstr string) {
	ip := net.ParseIP(ipstr)
	c := &mmdbformat{}
	err := db.Lookup(ip, c)
	if err != nil {
		ch <- getDefaultCountryData()
		return
	}

	countrydata := &models.DbCountry{
		CountryName: c.Country.Names["en"],
		CountryCode: c.Country.IsoCode,
		Continent:   c.Continent.Names["en"],
	}
	if c.Country.IsoCode == "US" {
		if len(c.Subdivisions) > 0 {
			countrydata.State = c.Subdivisions[0].IsoCode
		} else {
			countrydata.State = "Unknown"
		}
	} else {
		countrydata.State = "None"
	}
	ch <- countrydata
}
