package util

// util.go - Various convenience and utility functions

import "os"

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

func RemoveBytesFromBeginning(b []byte, num int) []byte {
	return b[num:]
}
