package conust

import (
	"math"
	"strconv"
	"strings"
)

const minBase = 2
const maxBase = 36

// FloatConverter can convert float64 numbers to strings in a base of your choice between 2 and 36, and back
// For integers you can use the go's strconv.FormatInt
type FloatConverter struct {
	base          int
	baseFloat     float64
	baseLog10     float64
	baseDecimals  int
	decDecimals   int
	builder       strings.Builder
	typicalLength int
}

// WithBaseDecimals sets how many fractional digits in the converter's base should be kept when converting to string
func (fc *FloatConverter) WithBaseDecimals(baseDecimals int) {
	fc.decDecimals = int(math.Floor(math.Abs(math.Log10(math.Pow(1.0/float64(fc.base), float64(fc.baseDecimals))))))
	fc.baseDecimals = baseDecimals
}

// WithDecimals sets a fractional precision high enough to ensure the given decimal precision
// When not set explicitly, the default is 3
func (fc *FloatConverter) WithDecimals(decDecimals int) {
	fc.baseDecimals = int(math.Ceil(math.Abs(float64(decDecimals) / math.Log10(float64(fc.base)))))
	fc.decDecimals = decDecimals
}

// WithTypicalLength sets the capaCity for the internal strings.Builder when ParseFloat is run
// When not set explicitly, the default is 10
// Choose a value that is large enough to accomodate (most of) your result strings but does not allocate needlessly large amounts of memory
func (fc *FloatConverter) WithTypicalLength(length int) {
	fc.typicalLength = length
}

// NewFloatConverter returns an initialized float64 base converter
// The base parameter must be between 2 and 36
func NewFloatConverter(base int) *FloatConverter {
	if base < minBase || base > maxBase {
		panic("invalid base parameter")
	}

	baseFloat := float64(base)
	fc := &FloatConverter{
		base:          base,
		baseFloat:     baseFloat,
		baseLog10:     math.Log10(baseFloat),
		typicalLength: 10,
	}
	fc.WithDecimals(3)
	return fc
}

// FormatFloat converts the number to a string in the specified base and with the fractional precision defined via WithDecimals or WithBaseDecimals
func (fc *FloatConverter) FormatFloat(input float64) string {
	fc.builder.Reset()
	fc.builder.Grow(fc.typicalLength)

	if input < 0 {
		input *= -1
		fc.builder.WriteByte(minusByte)
	}
	intPart, fracPart := math.Modf(input)
	fc.builder.WriteString(strconv.FormatInt(int64(intPart), fc.base))

	if fracPart != 0 {
		fc.builder.WriteByte(positiveIntegerTerminator)

		var currentDigit float64
		trailingZeros := 0
		for i := fc.baseDecimals; i > 0; i-- {
			fracPart *= fc.baseFloat
			if i == 1 {
				// round the last digit
				currentDigit = math.Round(fracPart)
			} else {
				currentDigit = math.Floor(fracPart)
				fracPart -= currentDigit
			}

			if currentDigit == 0 {
				trailingZeros++
				continue
			}

			if trailingZeros > 0 {
				for z := trailingZeros; z > 0; z-- {
					fc.builder.WriteByte(digit0)
				}
				trailingZeros = 0
			}

			fc.builder.WriteByte(intToDigit(int(currentDigit)))
		}
	}

	return fc.builder.String()
}

// TODO a string -> string conversion that rounds to decDecimals digits

// ParseFloat converts a numeric string to a floa64 number
func (fc *FloatConverter) ParseFloat(input string) (out float64, ok bool) {
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

func (fc *FloatConverter) isDigitLegal(digit int) bool {
	return digit >= 0 && digit < fc.base
}
