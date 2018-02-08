// Package entropy contains functions for checking the entropy of given data
package entropy

import (
	"bytes"
	"math"
	"strings"
)

// Define entropy thresholds over which a string is considered complex enough
// to be a potential key
const (
	Base64Threshold = 4.5
	HexThreshold    = 3.0
)

const (
	consider = 20 // When scanning for strings, only consider >= this value length
)

// CalculateShannon calculates the shannon entropy for a block of data
// - http://blog.dkbza.org/2007/05/scanning-data-for-entropy-anomalies.html
func CalculateShannon(data []byte) float64 {
	if len(data) == 0 {
		return 0.0
	}
	entropy := 0.0
	pX := 0.0
	for x := 0; x < 256; x++ {
		pX = float64(bytes.Count(data, []byte(string(x)))) / float64(len(data))
		if pX > 0 {
			entropy += -pX * math.Log2(pX)
		}
	}
	return entropy
}

// Check searches through a given block of data to attempt to identify high
// entropy blocks. Returns true and number of matching strings if found
func Check(b []byte) (bool, int) {

	found := [][]byte{}

	// Offset where we started reading the data - indexes from
	// 1 instead of zero otherwise we'll capture a spurious leading
	// byte into the slice
	// start := 1
	start := -1

	// Base64 strings
	for i, tok := range b {
		if !isBase64Byte(tok) || i+1 == len(b) {
			if i-start >= consider {
				s := b[start+1 : i]
				if e := CalculateShannon(s); e > Base64Threshold {
					found = append(found, s)
				}
			}
			start = i
		}
	}

	start = -1

	// Hex strings
	for i, tok := range b {
		if !isHexByte(tok) || i+1 == len(b) {
			if i-start >= consider {
				s := b[start+1 : i]
				if e := CalculateShannon(s); e > HexThreshold {
					found = append(found, s)
				}
			}
			start = i
		}
	}

	return len(found) == 0, len(found)
}

func isBase64Byte(b byte) bool {
	return strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=", string(b))
}

func isHexByte(b byte) bool {
	return strings.Contains("ABCDEFabcdef0123456789", string(b))
}
