package conust

import (
	"fmt"
	"strconv"
	"strings"
)

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

const signPositive byte = '9'
const signNegative byte = '0'
const zeroOutput = "5"

const builderInitialCap = 7

const trailing0 = "0"

// TODO are these decimal points ok?
const positiveIntegerTerminator byte = '.'
const negativeIntegerTerminator byte = '~'

// B10FromI32 encodes int32 into sortable string using decimal digits
func B10FromI32(i int32) (s string) {
	if i == 0 {
		return zeroOutput
	}
	return b10FromIntString(int32Preproc(i))
}

// B10ToI32 decodes string into int32
// Successfulness of the decoding is signalled by the second return value. A failure is possible when the encoded number is out of the range of the int32 type.
func B10ToI32(s string) (i int32, ok bool) {
	if s == zeroOutput {
		return 0, true
	}
	intPart, _ := decodeStrings(s, true, flipDigit10)
	result, err := strconv.ParseInt(intPart, 10, 32)
	return int32(result), (err == nil)
}

// B10FromI64 encodes int64 into sortable string using decimal digits
func B10FromI64(i int64) (s string) {
	if i == 0 {
		return zeroOutput
	}
	return b10FromIntString(int64Preproc(i))
}

// B10ToI64 decodes string into int32
// Successfulness of the decoding is signalled by the second return value. A failure is possible when the encoded number is out of the range of the int32 type.
func B10ToI64(s string) (i int64, ok bool) {
	if s == zeroOutput {
		return 0, true
	}
	intPart, _ := decodeStrings(s, true, flipDigit10)
	result, err := strconv.ParseInt(intPart, 10, 64)
	return result, (err == nil)
}

// B36FromI32 encodes int32 into sortable string using Base(36) digits
func B36FromI32(i int32) (s string) {
	return B36FromI64(int64(i))
}

// B36FromI64 encodes int64 into sortable string using Base(36) digits
func B36FromI64(i int64) (s string) {
	if i == 0 {
		return zeroOutput
	}
	var b strings.Builder
	b.Grow(builderInitialCap)
	var number string
	if i > 0 {
		b.WriteByte(signPositive)
		number = strconv.FormatInt(i, 36)
		intStringToB36(&b, true, number)
	} else {
		b.WriteByte(signNegative)
		number = strconv.FormatInt(i*-1, 36)
		intStringToB36(&b, false, number)
	}
	return b.String()
}

// B36ToI32 decodes string into int32
// Successfulness of the decoding is signalled by the second return value. A failure is possible when the encoded number is out of the range of the int32 type.
func B36ToI32(s string) (i int32, ok bool) {
	if s == zeroOutput {
		return 0, true
	}
	intPart, _ := decodeStrings(s, true, flipDigit36)
	result, err := strconv.ParseInt(intPart, 36, 32)
	return int32(result), (err == nil)
}

// B36ToI64 decodes string into int32
// Successfulness of the decoding is signalled by the second return value. A failure is possible when the encoded number is out of the range of the int32 type.
func B36ToI64(s string) (i int64, ok bool) {
	if s == zeroOutput {
		return 0, true
	}
	intPart, _ := decodeStrings(s, true, flipDigit36)
	result, err := strconv.ParseInt(intPart, 36, 64)
	fmt.Println("b36 decode raw", intPart, "error", err)
	return result, (err == nil)
}

func intStringToB10(b *strings.Builder, positive bool, number string) {
	if positive {
		b.WriteByte(intToDigit36(len(number)))
		b.WriteString(strings.TrimRight(number, trailing0))
	} else {
		b.WriteByte(intToReversedDigit36(len(number)))
		number = strings.TrimRight(number, trailing0)
		for j := 0; j < len(number); j++ {
			b.WriteByte(flipDigit10(number[j]))
		}
		b.WriteByte(negativeIntegerTerminator)
	}
}

func intStringToB36(b *strings.Builder, positive bool, number string) {
	if positive {
		b.WriteByte(intToDigit36(len(number)))
		b.WriteString(strings.TrimRight(number, trailing0))
	} else {
		b.WriteByte(intToReversedDigit36(len(number)))
		number = strings.TrimRight(number, trailing0)
		for j := 0; j < len(number); j++ {
			b.WriteByte(flipDigit36(number[j]))
		}
		b.WriteByte(negativeIntegerTerminator)
	}
}

func int32Preproc(i int32) (positive bool, number string) {
	if i > 0 {
		return true, fmt.Sprintf("%d", i)
	}
	return false, fmt.Sprintf("%d", i*-1)
}

func int64Preproc(i int64) (positive bool, absNumber string) {
	if i > 0 {
		return true, fmt.Sprintf("%d", i)
	}
	return false, fmt.Sprintf("%d", i*-1)
}

func b10FromIntString(positive bool, absNumber string) string {
	var b strings.Builder
	b.Grow(builderInitialCap)
	if positive {
		b.WriteByte(signPositive)
		intStringToB10(&b, true, absNumber)
	} else {
		b.WriteByte(signNegative)
		intStringToB10(&b, false, absNumber)
	}
	return b.String()
}

func flipDigit10(digit byte) byte {
	return intToReversedDigit10(digitToInt(digit))
}

func flipDigit36(digit byte) byte {
	return intToReversedDigit36(digitToInt(digit))
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

func intToDigit36(i int) byte {
	return digits36[i]
}

func intToReversedDigit10(i int) byte {
	return digits10Reversed[i]
}

func intToReversedDigit36(i int) byte {
	return digits36Reversed[i]
}

func decodeStrings(s string, intOnly bool, flipDigit func(byte) byte) (integral string, fractional string) {
	isPositive := (s[0] == signPositive)
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
