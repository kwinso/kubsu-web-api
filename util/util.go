package util

import (
	"crypto/rand"
	"encoding/hex"
	"math"
)

func Int32ToInt(input []int32) []int {
	output := make([]int, len(input))
	for i, v := range input {
		output[i] = int(v)
	}
	return output
}

// Uses cryptographically secure random number generator to generate username & password
func GenerateRandomCredentials() (string, string) {
	return randomBase16String(10), randomBase16String(10)
}

func randomBase16String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/2)))
	rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l] // strip 1 extra character we get from odd length results
}
