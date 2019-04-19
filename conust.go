// Package conust transforms numbers into string tokens for which the simple string comparison
// produces the same result as the numeric comparison of the original numbers would have.
// The input and the token representation can have different numeric bases in which case
// base conversion is done on the input. Base 36 is the highest possible base for both the input
// and the token.
//
// Transforming tokens back into numbers is also possible.
//
// The conversion adds to the length of the original numeric string, but at the same time it can
// save some space by omitting the trailing zeros of the integer part of the numbers as well as
// the leading zeros of the fractional part. Doing this results in the limitation that the integer
// part can not be longer than 35 digits and the fractional part can not have more than
// 35 leading zeros.
//
// You can make your numbers more compact by converting them to base 36. If you want to convert
// floating point numbers, you can use the FloatConverter of this package. For integers, there is
// the built in strconv.FormatInt function.

package conust

// [48 49 50 51 52 53 54 55 56 57 97 98 99 100 101 102 103 104 105 106 107 108 109 110 111 112 113 114 115 116 117 118 119 120 121 122]
var digits36 = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

var digits36Reversed = [...]byte{'z', 'y', 'x', 'w', 'v', 'u', 't', 's', 'r', 'q',
	'p', 'o', 'n', 'm', 'l', 'k', 'j', 'i', 'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a',
	'9', '8', '7', '6', '5', '4', '3', '2', '1', '0'}

const maxDigitValue = 35

const digit0 byte = '0'
const digit9 byte = '9'
const digitA byte = 'a'
const digitZ byte = 'z'
const minusByte byte = '-'
const plusByte byte = '+'

const signNegative byte = '4'
const zeroOutput = "5"
const signPositive byte = '6'

// LessThanAny is a string which is less than any encoded value
const LessThanAny = "3"

// GreaterThanAny is a string which is greater than any encoded value
const GreaterThanAny = "7"

const zeroInput = "0"

const positiveIntegerTerminator byte = '.'
const negativeIntegerTerminator byte = '~'

// Codec can transform strings to and from the Conust format
type Codec interface {
	Encode(in string) (out string, ok bool)
	Decode(in string) (out string, ok bool)
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

func flipDigit(digit byte) byte {
	return intToReversedDigit(digitToInt(digit))
}

func intToReversedDigit(i int) byte {
	return digits36Reversed[i]
}

func intToDigit(i int) byte {
	return digits36[i]
}
