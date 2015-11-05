package util

func ReadTillNul(b []byte) string {
	var s []byte
	for i, _ := range b {
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
