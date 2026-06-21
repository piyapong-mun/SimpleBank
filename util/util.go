package util

// For random string generation
import (
	"math/rand"
	"time"
)

// RandomString generates a random string of a given length
func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// RandomNumber generates a random number between min and max
func RandomNumber(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min+1) + min
}

func RandomCurrency() string {
	var currencies = []string{"USD", "EUR", "THB"}
	rand.Seed(time.Now().UnixNano())
	return currencies[rand.Intn(len(currencies))]
}
