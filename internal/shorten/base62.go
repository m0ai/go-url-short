package shorten

import (
	"fmt"
	"math"
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

// ConvertRadix10 converts a integer from an base62 string
func ConvertRadix10(shortKey string) (int64, error) {
	var id int64 = 0

	shortKey = reverse(shortKey)
	for i, c := range shortKey {
		index := strings.IndexRune(base62charset, c)
		if index == -1 {
			return 0, ErrInvalidID{fmt.Errorf("invalid character: %c", c)}
		}
		id += int64(index) * int64(math.Pow(float64(base62keyLength), float64(i)))

	}

	return id, nil
}

// GenerateBase62 converts a string from an integer
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
