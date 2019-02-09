package conust

import (
	"fmt"
	"strconv"
	"strings"
)

// Encoder turns numbers or decimal strings to Conust strings
type Encoder interface {
	FromString(string) (string, bool)
	FromInt32(int32) string
	FromInt64(int64) string
	FromFloat32(float32) string
	FromFloat64(float64) string
}

// Decoder turns Conust strings back to numbers or decimal strings
type Decoder interface {
	ToString(string) (string, bool)
	ToInt32(string) (int32, bool)
	ToInt64(string) (int64, bool)
	ToFloat32(string) (float32, bool)
	ToFloat64(string) (float64, bool)
}

// not used, it is a subset of digits36
// var digits10 = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

var digits10Reversed = [...]byte{'9', '8', '7', '6', '5', '4', '3', '2', '1', '0'}

// [48 49 50 51 52 53 54 55 56 57 97 98 99 100 101 102 103 104 105 106 107 108 109 110 111 112 113 114 115 116 117 118 119 120 121 122]
var digits36 = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

var digits36Reversed = [...]byte{'z', 'y', 'x', 'w', 'v', 'u', 't', 's', 'r', 'q',
	'p', 'o', 'n', 'm', 'l', 'k', 'j', 'i', 'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a',
	'9', '8', '7', '6', '5', '4', '3', '2', '1', '0'}

const digit0 byte = '0'
const digit9 byte = '9'
const digitA byte = 'a'
const digitZ byte = 'z'
const minusByte byte = '-'
const plusByte byte = '+'

const signNegative byte = '4'
const zeroOutput = "5"
const signPositive byte = '6'

const builderInitialCap = 7

const trailing0 = "0"

const positiveIntegerTerminator byte = '.'
const negativeIntegerTerminator byte = '~'

func int32Preproc(i int32) (positive bool, absNumber string) {
	if i < 0 {
		return false, fmt.Sprintf("%d", i)[1:]
	}
	return true, fmt.Sprintf("%d", i)
}

func int64Preproc(i int64) (positive bool, absNumber string) {
	if i < 0 {
		return false, fmt.Sprintf("%d", i)[1:]
	}
	return true, fmt.Sprintf("%d", i)
}

func float32Preproc(f float32, precision int) (positive bool, absNumber string) {
	if f < 0 {
		return false, strconv.FormatFloat(float64(f), 'f', precision, 32)[1:]
	}
	return true, strconv.FormatFloat(float64(f), 'f', precision, 32)
}

func float64Preproc(f float64, precision int) (positive bool, absNumber string) {
	if f < 0 {
		return false, strconv.FormatFloat(f, 'f', precision, 64)[1:]
	}
	return true, strconv.FormatFloat(f, 'f', precision, 64)
}

func digitToInt(digit byte) int {
	if digit < digitA {
		return int(digit - digit0)
	}
	return 10 + int(digit-digitA)
}

func reversedDigitToInt(digit byte) int {
	if digit < digitA {
		return 26 + int(digit9-digit)
	}
	return int(digitZ - digit)
}

func intToDigit(i int) byte {
	return digits36[i]
}

func intToReversedDigit10(i int) byte {
	return digits10Reversed[i]
}

func intToReversedDigit36(i int) byte {
	return digits36Reversed[i]
}

func flipDigit10(digit byte) byte {
	return intToReversedDigit10(digitToInt(digit))
}

func flipDigit36(digit byte) byte {
	return intToReversedDigit36(digitToInt(digit))
}

func decodeStrings(s string, intOnly bool, flipDigit func(byte) byte) (integral string, fractional string) {
	isPositive := (s[0] > zeroOutput[0])
	var terminator byte
	var intLength int

	if isPositive {
		terminator = positiveIntegerTerminator
		intLength = digitToInt(s[1])
	} else {
		terminator = negativeIntegerTerminator
		intLength = reversedDigitToInt(s[1])
	}

	var b strings.Builder
	b.Grow(intLength + 1)

	if isPositive {
		// b.WriteByte('+')
	} else {
		b.WriteByte('-')
	}

	terminatorPos := strings.IndexByte(s, terminator)
	if terminatorPos > 1 {
		if isPositive {
			b.WriteString(s[2:terminatorPos])
		} else {
			for i := 2; i < terminatorPos; i++ {
				b.WriteByte(flipDigit(s[i]))
			}
		}
	} else {
		if isPositive {
			b.WriteString(s[2:])
		} else {
			for i := 2; i < len(s); i++ {
				b.WriteByte(flipDigit(s[i]))
			}
		}
	}
	if bLength := b.Len(); bLength <= intLength {
		if !isPositive {
			bLength--
		}
		for i := bLength; i < intLength; i++ {
			b.WriteByte('0')
		}
	}
	if intOnly || terminatorPos < 0 || terminatorPos == len(s)-1 {
		return b.String(), ""
	}

	integral = b.String()
	b.Reset()
	leadingZeroes := reversedDigitToInt(s[terminatorPos+1])
	for i := 0; i < leadingZeroes; i++ {
		b.WriteByte('0')
	}
	b.WriteString(s[terminatorPos+2:])
	return integral, b.String()
}
