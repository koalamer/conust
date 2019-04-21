package conust

import (
	"math"
	"strconv"
)

const minBase = 2
const maxBase = 36

// FloatConverter can convert float64 numbers to strings in a base of your choice between 2 and 36, and back
// For integers you can use the go's strconv.FormatInt
type FloatConverter interface {
	// WithBaseDecimals sets how many fractional digits in the converter's base should be kept when converting to string
	WithBaseDecimals(baseDecimals int) *floatConverter
	// WithDecimals sets a fractional precision high enough to ensure the given decimal precision
	WithDecimals(decDecimals int) *floatConverter
	// NewFloatConverter returns an initialized float64 base converter
	FormatFloat(input float64) string
	// ParseFloat converts a numeric string to a floa64 number
	ParseFloat(input string) (out float64, ok bool)
}

type floatConverter struct {
	base         int
	baseFloat    float64
	baseLog10    float64
	baseDecimals int
	decDecimals  int
	precision    float64
	buffer       [80]byte
}

// WithBaseDecimals sets how many fractional digits in the converter's base should be kept when converting to string
func (fc *floatConverter) WithBaseDecimals(baseDecimals int) *floatConverter {
	fc.decDecimals = int(math.Floor(math.Log10(math.Pow(fc.baseFloat, float64(fc.baseDecimals)))))
	fc.baseDecimals = baseDecimals
	fc.precision = math.Pow10(-fc.decDecimals)
	return fc
}

// WithDecimals sets a fractional precision high enough to ensure the given decimal precision
// When not set explicitly, the default is 3
func (fc *floatConverter) WithDecimals(decDecimals int) *floatConverter {
	fc.baseDecimals = int(math.Ceil(math.Abs(float64(decDecimals) / math.Log10(fc.baseFloat))))
	fc.decDecimals = decDecimals
	fc.precision = math.Pow10(-fc.decDecimals)
	return fc
}

// NewFloatConverter returns an initialized float64 base converter
// The base parameter must be between 2 and 36
func NewFloatConverter(base int) FloatConverter {
	if base < minBase || base > maxBase {
		panic("invalid base parameter")
	}

	baseFloat := float64(base)
	fc := &floatConverter{
		base:      base,
		baseFloat: baseFloat,
		baseLog10: math.Log10(baseFloat),
	}
	fc.WithDecimals(3)
	return fc
}

// FormatFloat converts the number to a string in the specified base and with the fractional precision defined via WithDecimals or WithBaseDecimals
func (fc *floatConverter) FormatFloat(input float64) string {
	input = fc.roundToPrecision(input)

	start := 1 // in case of an overflow, first byte might be needed
	cursor := 1

	isNegative := input < 0
	if isNegative {
		input *= -1
		fc.buffer[cursor] = minusByte
		cursor++
	}

	intPart, fracPart := math.Modf(input)

	intString := strconv.FormatInt(int64(intPart), fc.base)
	intStringDest := fc.buffer[cursor : cursor+len(intString)]
	copy(intStringDest, intString)
	cursor += len(intString)

	if fracPart != 0 {
		lastDigitOverflows := false
		fc.buffer[cursor] = positiveIntegerTerminator
		cursor++
		fracPartStart := cursor

		var currentDigit float64
		trailingZeros := 0
		for i := fc.baseDecimals; i > 0; i-- {
			fracPart *= fc.baseFloat
			currentDigit = math.Floor(fracPart)
			fracPart -= currentDigit

			if i == 1 {
				// round the last digit to make decDecimals guarantee work without using additional digits
				remainder := fracPart
				if isNegative {
					remainder *= -1
				}
				remainder = math.Abs(math.Round(remainder))
				if remainder > 0 {
					if currentDigit < fc.baseFloat-1 {
						currentDigit++
					} else {
						lastDigitOverflows = true
					}

				}
			}

			if currentDigit == 0 {
				trailingZeros++
				continue
			}

			if trailingZeros > 0 {
				for z := trailingZeros; z > 0; z-- {
					fc.buffer[cursor] = digit0
					cursor++
				}
				trailingZeros = 0
			}

			fc.buffer[cursor] = intToDigit(int(currentDigit))
			cursor++
		}

		if cursor == fracPartStart {
			// no fractional part was written, only the decimal separator
			cursor--
		} else if lastDigitOverflows {
			start, cursor = fc.applyOverflow(start, cursor)
		}
	}

	return string(fc.buffer[start:cursor])
}

// TODO a string -> string conversion that rounds to decDecimals digits

// ParseFloat converts a numeric string to a floa64 number
func (fc *floatConverter) ParseFloat(input string) (out float64, ok bool) {
	isPositive := true
	length := len(input)
	var cursor int

	if input[0] == minusByte {
		isPositive = false
		cursor++
	} else if input[0] == plusByte {
		cursor++
	}

	for ; cursor < length; cursor++ {
		c := input[cursor]
		if c == positiveIntegerTerminator {
			cursor++
			break
		}
		digit := digitToInt(c)
		if !fc.isDigitLegal(digit) {
			return 0, false
		}
		out = out*fc.baseFloat + float64(digit)
	}

	var fraction float64 = 1
	for ; cursor < length; cursor++ {
		fraction /= fc.baseFloat
		digit := digitToInt(input[cursor])
		if !fc.isDigitLegal(digit) {
			return 0, false
		}
		if digit == 0 {
			continue
		}
		out += float64(digit) * fraction
	}

	if !isPositive {
		out *= -1
	}
	return out, true
}

func (fc *floatConverter) isDigitLegal(digit int) bool {
	return digit >= 0 && digit < fc.base
}

func (fc *floatConverter) applyOverflow(start int, end int) (newStart int, newEnd int) {
	newStart = start
	newEnd = end - 1
	maxDigit := intToDigit(fc.base - 1)

	// fraction part
	for i := newEnd; i >= start; i-- {
		switch fc.buffer[i] {
		case maxDigit:
			newEnd = i
		case positiveIntegerTerminator:
			newEnd = i
			i = -1 // aka stop the for loop
		default:
			fc.buffer[i] = intToDigit(digitToInt(fc.buffer[i]) + 1)
			return
		}
	}

	// integer part
	isNegative := fc.buffer[start] == minusByte
	var intStart int
	if isNegative {
		intStart = start + 1
	} else {
		intStart = start
	}
	for i := newEnd - 1; i >= intStart; i-- {
		switch fc.buffer[i] {
		case maxDigit:
			fc.buffer[i] = digit0
		default:
			fc.buffer[i] = intToDigit(digitToInt(fc.buffer[i]) + 1)
			return
		}
	}

	// reached the start
	newStart = start - 1
	if isNegative {
		fc.buffer[newStart] = minusByte
		fc.buffer[start] = digit1
	} else {
		fc.buffer[newStart] = digit1
	}
	return
}

func (fc *floatConverter) roundToPrecision(input float64) float64 {
	precisionFraction := input / fc.precision
	precisionFraction -= math.Round(precisionFraction)
	return input - fc.precision*precisionFraction
}
