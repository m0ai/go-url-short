package shorten

import (
	"fmt"
	"strings"
)

type ErrInvalidID struct {
	Err error
}

func (e ErrInvalidID) Error() string {
	return "invalid ID"
}

const base62charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base62keyLength = len(base62charset)

// GenerateBase62 generates a base62 string from an integer
func ConvertRadix62(id int64) string {
	result := make([]string, 0)
	var quotient, remainder int64 = id, 0

	for quotient > 0 {
		quotient, remainder = divmod(quotient, int64(base62keyLength))
		ch := base62charset[remainder]
		fmt.Println(remainder, ch)
		result = append(result, string(ch))
	}

	return reverse(strings.Join(result, ""))
}

func divmod(numerator, denominator int64) (quotient, remainder int64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
