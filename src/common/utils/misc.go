package utils

import (
	"fmt"
	"crypto/rand"
)

func NewRandLenChars(length int) string {
	r, _ := NewRandLenCustomChars(length, stdChars)
	return r
}

var stdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

func NewRandLenCustomChars(length int, chars []byte) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("[newRandLenChars] Wrong Length!")
	}
	clen := len(chars)
	if clen < 2 || clen > 256 {
		return "", fmt.Errorf("[newRandLenChars] Wrong Charset Length!")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4))
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			return "", fmt.Errorf("[newRandLenChars] Error Reading Random Bytes! error=[%v]", err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue
			}
			b[i] = chars[c%clen]
			i++
			if i == length {
				return string(b), nil
			}
		}
	}
}