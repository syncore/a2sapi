package db

import (
	"database/sql"
	"fmt"
	"net"

	_ "github.com/mattn/go-sqlite3"
)

type Country struct {
	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`
	Region      string `json:"region"`
}

func OpenCountryDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "countries.sqlite")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetCountryInfo(db *sql.DB, ip string) (*Country, error) {
	ctry := GetDefaultCountryData()

	i := net.ParseIP(ip).To4()
	ipint := uint32(i[0])<<24 | uint32(i[1])<<16 | uint32(i[2])<<8 | uint32(i[3])
	query := `SELECT country_location_data.continent_name,
	country_location_data.country_iso_code,country_location_data.country_name FROM
	ip_range_data INNER JOIN country_location_data ON ip_range_data.geoname_id =
	country_location_data.geoname_id WHERE ip_range_data.range_begin <=?
	AND ip_range_data.range_end >=? LIMIT 1`

	rows, err := db.Query(query, ipint, ipint)
	if err != nil {
		fmt.Printf("DB error detected in GetCountryInfo: %s\n", err)
		return ctry, err
	}
	defer rows.Close()
	for rows.Next() {
		var continent string
		var countrycode string
		var countryname string

		if err := rows.Scan(&continent, &countrycode, &countryname); err != nil {
			fmt.Printf("DB error detected in GetCountryInfo: %s\n", err)
			return ctry, err
		}
		ctry.CountryName = countryname
		ctry.CountryCode = countrycode
		ctry.Region = continent
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("DB error detected in GetCountryInfo: %s\n", err)
		return GetDefaultCountryData(), err
	}
	return ctry, nil
}

func GetDefaultCountryData() *Country {
	return &Country{
		CountryName: "Unknown",
		CountryCode: "Unknown",
		Region:      "Unknown",
	}
}
