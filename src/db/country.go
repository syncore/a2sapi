package db

// country.go - Country geolocation database lookup.

import (
	"a2sapi/src/constants"
	"a2sapi/src/logger"
	"a2sapi/src/models"
	"fmt"
	"net"
	"runtime"

	"github.com/oschwald/maxminddb-golang"
)

// CDB represents a database containing geolocation information.
type CDB struct {
	db *maxminddb.Reader
}

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

func getDefaultCountryData() models.DbCountry {
	return models.DbCountry{
		CountryName: "Unknown",
		CountryCode: "Unknown",
		Continent:   "Unknown",
		State:       "Unknown",
	}
}

// OpenCountryDB opens the country lookup database for reading. The caller of
// this function will be responsinble for calling .Close().
func OpenCountryDB() (*CDB, error) {
	// Note: the caller of this function needs to handle .Close()
	conn, err := maxminddb.Open(constants.CountryDbFilePath)
	if err != nil {
		dir := "build/nix"
		if runtime.GOOS == "windows" {
			dir = "build\\win"
		}
		logger.LogAppError(err)
		panic(
			fmt.Sprintf(
				`Unable to open country database! Use the get_countrydb script in the %s
directory to get the country DB file, or download it from:
http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz and
extract the "GeoLite2-City.mmdb" file into a directory called "db" in the same
directory as the a2sapi executable. Error: %s`, dir, err))
	}
	return &CDB{db: conn}, nil
}

// Close closes the country geolocation database.
func (cdb *CDB) Close() {
	err := cdb.db.Close()
	if err != nil {
		logger.LogAppErrorf("Error closing country database: %s", err)
	}
}

// GetCountryInfo attempts to retrieve the country information for a given IP,
// returning the result as a country model object over the corresponding result channel.
func (cdb *CDB) GetCountryInfo(ch chan<- models.DbCountry, ipstr string) {
	ip := net.ParseIP(ipstr)
	c := &mmdbformat{}
	err := cdb.db.Lookup(ip, c)
	if err != nil {
		ch <- getDefaultCountryData()
		return
	}
	if c.Country.Names["en"] == "" || c.Country.IsoCode == "" {
		ch <- getDefaultCountryData()
		return
	}

	countrydata := models.DbCountry{
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
