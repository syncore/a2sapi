package db

// country.go - Country geolocation database lookup.

import (
	"net"
	"steamtest/src/util"

	"github.com/oschwald/maxminddb-golang"
)

const mmDbFile = "GeoLite2-City.mmdb"

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

// The Country struct is for the JSON representation.
type Country struct {
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	Continent   string `json:"region"`
	State       string `json:"state"`
}

func getDefaultCountryData() *Country {
	return &Country{
		CountryName: "Unknown",
		CountryCode: "Unknown",
		Continent:   "Unknown",
		State:       "Unknown",
	}
}

func OpenCountryDB() (*maxminddb.Reader, error) {
	// Note: the caller of this function needs to handle db.Close()
	db, err := maxminddb.Open(mmDbFile)
	if err != nil {
		return nil, util.LogAppError(err.Error())
	}
	return db, nil
}

func GetCountryInfo(ch chan<- *Country, db *maxminddb.Reader, ipstr string) {
	ip := net.ParseIP(ipstr)
	c := &mmdbformat{}
	err := db.Lookup(ip, c)
	if err != nil {
		ch <- getDefaultCountryData()
		return
	}

	countrydata := &Country{
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
