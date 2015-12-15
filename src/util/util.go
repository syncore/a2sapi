package util

// util.go - Various convenience and utility functions

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func DirExists(dir string) bool {
	f, err := os.Stat(dir)
	return err == nil && f.IsDir()
}

func CreateDirectory(dirname string) error {
	if DirExists(dirname) {
		return nil
	}
	if err := os.Mkdir(dirname, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func WriteJSONConfig(cfg interface{}, directory, fullpath string) error {
	if err := CreateDirectory(directory); err != nil {
		fmt.Printf("Couldn't create '%s' dir: %s\n", directory, err)
		return err
	}
	f, err := os.Create(fullpath)
	if err != nil {
		fmt.Printf("Couldn't create '%s' file: %s\n", fullpath, err)
		return err
	}
	defer f.Close()
	c, err := json.Marshal(cfg)
	if err != nil {
		fmt.Printf("Error marshaling JSON for '%s' file: %s\n", fullpath, err)
		return err
	}
	w := bufio.NewWriter(f)
	_, err = w.Write(c)
	if err != nil {
		fmt.Printf("Error writing contents of config file to disk: %s\n", err)
		return err
	}
	f.Sync()
	w.Flush()
	fmt.Printf("Successfuly wrote config file to %s\n\n", fullpath)
	return nil
}

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
