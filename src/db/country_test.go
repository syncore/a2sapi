package db

import (
	"os"
	"steamtest/src/config"
	"steamtest/src/constants"
	"steamtest/src/models"
	"strings"
	"testing"
)

func init() {
	// need base directory for config and other files
	err := os.Chdir("../..")
	if err != nil {
		panic("Unable to change directory for tests")
	}
	doCleanup()
	// use testing configuration
	config.CreateTestConfig()
	constants.IsTest = true
}

func TestOpenCountryDB(t *testing.T) {
	db, err := OpenCountryDB()
	// Will panic anyway
	if err != nil {
		t.Fatalf("Error opening country database: %s", err)
	}
	defer db.Close()
}

func TestGetCountryInfo(t *testing.T) {
	cdb, err := OpenCountryDB()
	if err != nil {
		t.Fatalf("Error opening country database: %s", err)
	}
	defer cdb.Close()
	c := make(chan *models.DbCountry, 1)
	ip := "192.211.62.11"
	cinfo := &models.DbCountry{}
	go GetCountryInfo(c, cdb, ip)
	cinfo = <-c
	if !strings.EqualFold(cinfo.CountryCode, "US") {
		t.Fatalf("Expected country code to be US for IP: %s, got: %s",
			ip, cinfo.CountryCode)
	}
	ip = "89.20.244.197"
	cinfo = &models.DbCountry{}
	go GetCountryInfo(c, cdb, ip)
	cinfo = <-c
	if !strings.EqualFold(cinfo.CountryCode, "NO") {
		t.Fatalf("Expected country code to be NO for IP: %s, got: %s",
			ip, cinfo.CountryCode)
	}
}
