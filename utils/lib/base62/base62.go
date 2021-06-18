package base62

import (
	"errors"
	"strings"
)

var (
	//Base ...
	// CharacterSet consists of 62 characters [0-9][A-Z][a-z].
	Base int64 = 62
	//CharacterSet ...
	CharacterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// Encode returns a base62 representation as
// string of the given integer number.
func Encode(num int64) string {
	b := make([]byte, 0)
	for num > 0 {
		r := int(num % Base)
		num /= Base
		b = append([]byte{CharacterSet[r]}, b...)
	}
	return string(b)
}

// Decode returns a integer number of a base62 encoded string.
func Decode(s string) (int64, error) {
	var r int64
	var power int
	for i, v := range s {
		power = len(s) - (i + 1)
		pos := strings.IndexRune(CharacterSet, v)
		if pos == -1 {
			return int64(pos), errors.New("invalid character: " + string(v))
		}
		r += int64(pos) * pow(Base, int64(power))
	}
	return r, nil
}

func pow(base int64, exponent int64) int64 {
	if exponent != 0 {
		return (base * pow(base, exponent-1))
	} else {
		return 1
	}
}
