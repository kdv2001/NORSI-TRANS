package converter

import "fmt"

func power(base int, degree int) int {
	result := 1
	for i := 0; i < degree; i++ {
		result *= base
	}
	return result
}

type HexConverter interface {
	Convert(hexString string) (decString string)
}

type hexConverter struct{}

func NewHexConverter() HexConverter {
	return &hexConverter{}
}

func (hc *hexConverter) Convert(hexString string) (decString string) {
	hexMap := map[byte]int{
		'0': 0,
		'1': 1,
		'2': 2,
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'a': 10,
		'b': 11,
		'c': 12,
		'd': 13,
		'e': 14,
		'f': 15,
		'A': 10,
		'B': 11,
		'C': 12,
		'D': 13,
		'E': 14,
		'F': 15,
	}

	hexLen := len(hexString)
	var decNum int
	for i, val := range hexString {
		hexDigit := hexMap[byte(val)]
		decNum += hexDigit * power(16, hexLen-i-1)
	}

	return fmt.Sprintf("%d", decNum)
}
