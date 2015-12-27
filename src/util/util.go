package util

// util.go - Various convenience and utility functions

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// FileExists returns true if the file of name exists; otherwise returns false.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// DirExists returns true if the directory dir exists; otherwise returns false.
func DirExists(dir string) bool {
	f, err := os.Stat(dir)
	return err == nil && f.IsDir()
}

// CreateDirectory attempts to create a directory with given dirname if it does
// not already exist. If it already exists, then the function does nothing.
func CreateDirectory(dirname string) error {
	if DirExists(dirname) {
		return nil
	}
	if err := os.Mkdir(dirname, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// CreateEmptyFile wraps os.Create and creates and empty file with at fullpath if
// it does not exist. If the file at fullpath exists, and allowOverwrite is true
// then it will overwrite the existing file.
func CreateEmptyFile(fullpath string, allowOverwrite bool) error {
	if FileExists(fullpath) && !allowOverwrite {
		return fmt.Errorf("File %s exists and allowOverwrite is false", fullpath)
	}
	f, err := os.Create(fullpath)
	if err != nil {
		return fmt.Errorf("Couldn't create empty file '%s': %s\n", fullpath, err)
	}
	defer f.Close()
	return nil
}

// CreateByteFile takes a slice of bytes & writes its contents to disk at
// fullpath using an io.Writer if the file does not already exist. If the
// file exists and allowOverwrite is true then it will overwrite the existing file.
func CreateByteFile(data []byte, fullpath string, allowOverwrite bool) error {
	if FileExists(fullpath) && !allowOverwrite {
		return fmt.Errorf("File %s exists and allowOverwrite is false", fullpath)
	}
	f, err := os.Create(fullpath)
	if err != nil {
		return fmt.Errorf("Couldn't create '%s' file: %s\n", fullpath, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("Error writing file to disk: %s\n", err)
	}
	f.Sync()
	w.Flush()
	return nil
}

// WriteJSONConfig takes an empty interface cfg, which is meant to be a
// *config.Config struct or a *GameList and writes it as JSON to disk at
// fullpath, overwriting any existing file.
func WriteJSONConfig(cfg interface{}, directory, fullpath string) error {
	if err := CreateDirectory(directory); err != nil {
		return fmt.Errorf("Couldn't create '%s' dir: %s\n", directory, err)
	}
	c, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Error marshaling JSON for '%s' file: %s\n", fullpath, err)
	}
	err = CreateByteFile(c, fullpath, true)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	fmt.Printf("Successfuly wrote config file to %s\n\n", fullpath)
	return nil
}

// ReadTillNul takes a slice of bytes b and reads up until it encounters the
// first null terminator, returning the bytes read up until that point as a string.
func ReadTillNul(b []byte) string {
	var s []byte
	for i := range b {
		if b[i] != '\x00' {
			s = append(s, b[i])
		} else {
			break
		}
	}
	return string(s)
}
