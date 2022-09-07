package faker_util

import (
	"fmt"
	"math/rand"
	"time"
)

// RGBColor RBG Color Type
type RGBColor struct {
	Red   int
	Green int
	Blue  int
}

// GetHex Converts a decimal number to hex representations
func getHex(num int) string {
	hex := fmt.Sprintf("%x", num)
	if len(hex) == 1 {
		hex = "0" + hex
	}
	return hex
}

// GetRandomColorInRgb Returns a random RGBColor
func GetRandomColorInRgb() RGBColor {
	rand.Seed(time.Now().UnixNano())
	Red := rand.Intn(255)
	Green := rand.Intn(255)
	blue := rand.Intn(255)
	c := RGBColor{Red, Green, blue}
	return c
}

// GetRandomColorInHex returns a random color in HEX format
func GetRandomColorInHex() string {
	color := GetRandomColorInRgb()
	hex := "#" + getHex(color.Red) + getHex(color.Green) + getHex(color.Blue)
	return hex
}
