package main

import (
	"NORSI-TRANS/converter"
	"fmt"
)

func main() {
	hc := converter.NewHexConverter()
	hexString := "123456789abcdef"
	decString := hc.Convert(hexString)
	fmt.Printf("Hex: %s, Dec: %s\n", hexString, decString)
}
