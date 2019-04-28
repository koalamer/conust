package conust

import (
	"fmt"
	"strings"
)

type codec struct {
	// internal analysis subject
	input string

	// internal anaysis state
	length               int
	cursor               int
	ok                   bool
	empty                bool
	zero                 bool
	positive             bool
	intSignificantFrom   int
	intSignificantTo     int
	intTo                int
	fracLeadingZeroCount int
	fracSignificantFrom  int
	fracSignificantTo    int
	magnitude            int
	magnitudePositive    bool

	// output builder
	builder strings.Builder
}

// NewCodec returns a codec
func NewCodec() Codec {
	return &codec{}
}

//Encode turns the input into the alphanumerically sortable Conust string
func (c *codec) Encode(s string) (out string, ok bool) {
	c.input = s
	c.AnalyzeInput()

	if !c.ok {
		c.input = ""
		return "", false
	}
	if c.empty {
		c.input = ""
		return "", true
	}
	if c.zero {
		c.input = ""
		return zeroOutput, true
	}

	intNonZeroLength := c.intSignificantTo - c.intSignificantFrom
	hasFraction := c.fracSignificantTo > c.fracSignificantFrom

	var numberNonZeroFrom int
	var numberNonZeroTo int
	var numberLength int
	if intNonZeroLength > 0 {
		// int size
		c.magnitudePositive = true
		c.magnitude = c.intTo - c.intSignificantFrom
		numberNonZeroFrom = c.intSignificantFrom
		if hasFraction {
			numberNonZeroTo = c.fracSignificantTo
			// decimal point won't be printed
			numberLength = numberNonZeroTo - numberNonZeroFrom - 1
		} else {
			numberNonZeroTo = c.intSignificantTo
			numberLength = numberNonZeroTo - numberNonZeroFrom
		}
	} else {
		// leading far zeros
		c.magnitudePositive = false
		c.magnitude = c.fracLeadingZeroCount
		numberNonZeroFrom = c.fracSignificantFrom
		numberNonZeroTo = c.fracSignificantTo
		numberLength = numberNonZeroTo - numberNonZeroFrom
	}
	magnitudeDigitCount := 0
	for i := c.magnitude; i > 0; i -= maxMagnitudeDigitValue {
		magnitudeDigitCount++
	}
	if magnitudeDigitCount == 0 {
		magnitudeDigitCount++
	}

	// determine needed allocation size
	outLength := 1 + magnitudeDigitCount + numberLength // sign + magnitude + length
	if !c.positive {
		// for the negative terminator
		outLength++
	}

	/*
		// abort if storage is exceeded when using static buffer
		if outLength > len(c.buffer) {
			c.input = ""
			return "", false
		}
	*/

	// build encoded string
	c.builder.Reset()
	c.builder.Grow(outLength)

	// write sign byte
	c.builder.WriteByte(c.encodeSign())

	// write magnitude bytes
	m := c.magnitude
	for i := magnitudeDigitCount; i > 0; i-- {
		if m > maxMagnitudeDigitValue {
			c.builder.WriteByte(c.encodeMagnitude(maxDigitValue))
		} else {
			c.builder.WriteByte(c.encodeMagnitude(m))
		}
		m -= maxMagnitudeDigitValue
	}

	if intNonZeroLength > 0 && hasFraction {
		if c.positive {
			// write integer part
			for i := numberNonZeroFrom; i < c.intTo; i++ {
				c.builder.WriteByte(c.input[i])
			}
			// write fraction part
			for i := c.intTo + 1; i < numberNonZeroTo; i++ {
				c.builder.WriteByte(c.input[i])
			}
		} else {
			// write integer part
			for i := numberNonZeroFrom; i < c.intTo; i++ {
				c.builder.WriteByte(flipDigit(c.input[i]))
			}
			// write fraction part
			for i := c.intTo + 1; i < numberNonZeroTo; i++ {
				c.builder.WriteByte(flipDigit(c.input[i]))
			}
			// negative terminator
			c.builder.WriteByte(negativeNumberTerminator)
		}
	} else {
		if c.positive {
			// non zero part
			for i := numberNonZeroFrom; i < numberNonZeroTo; i++ {
				c.builder.WriteByte(c.input[i])
			}
		} else {
			// non zero part
			for i := numberNonZeroFrom; i < numberNonZeroTo; i++ {
				c.builder.WriteByte(flipDigit(c.input[i]))
			}
			// negative terminator
			c.builder.WriteByte(negativeNumberTerminator)
		}
	}
	out = c.builder.String()
	ok = true
	c.input = ""

	// TODO remove debug
	if outLength != len(out) {
		panic(fmt.Sprintf("outLength error: for %v -> %v expected %d, got %d", s, out, outLength, len(out)))
	}

	return
}

// AnalyzeInput produces the correct internal state for the encoding step
func (c *codec) AnalyzeInput() {
	/* defer func() {
		fmt.Printf("Input: %q, ok: %v, empty: %-5v, zero: %-5v, positive: %-5v\n", s, ok, c.isEmpty, c.isZero, c.isPositive)
		fmt.Printf("  int nz start: %-3d, nz end: %-3d, end: %-3d\n", c.intNonZeroFrom, c.intNonZeroTo, c.intTo)
		fmt.Printf("  frac leading z count: %-3d, nz start: %-3d, nz end: %-3d\n", c.fracLeadingZeroCount, c.fracNonZeroFrom, c.fracNonZeroTo)
		fmt.Printf("  base %d -> %d, precision %d -> %d, conversion needed: %v\n", c.sourceBase, c.newBase, c.sourceBasePrecision, c.newBasePrecision, c.baseConversionNeeded)
	}() */
	c.cursor = 0
	c.ok = false
	c.length = len(c.input)

	// empty input results in empty but ok output
	if c.length == 0 {
		c.ok = true
		c.empty = true
		return
	}

	c.empty = false
	c.zero = true

	if c.checkSign() {
		c.cursor++
	}
	// a sign only is bad input
	if !c.cursorCanRead() {
		return
	}

	// skip leading zeroes
	c.skipZeroes()

	// if there were only zeroes, the result is zeroOutput (c.zero = true)
	if !c.cursorCanRead() {
		c.ok = true
		return
	}

	// determine integer part bounds
	if !c.getIntPartBounds() {
		return
	}

	c.resetFractionParams()

	// if no fraction present, end processing
	if !c.cursorCanRead() {
		c.ok = true
		return
	}
	// if the last non digit chaccter is not the decimal separator, that's an error
	if c.input[c.cursor] != decimalPoint {
		return
	}
	// skip over decimal separator
	c.cursor++

	// process fraction part
	c.ok = c.getFractionPartBounds()
	return
}

// Decode turns a Conust string back into its normal representation
func (c *codec) Decode(s string) (out string, ok bool) {
	c.input = s
	c.AnalyzeToken(s)

	if !c.ok {
		c.input = ""
		return "", false
	}
	if c.empty {
		c.input = ""
		return "", true
	}
	if c.zero {
		c.input = ""
		return zeroInput, true
	}

	// calculate output length
	outLength := c.intTo - c.intSignificantFrom
	if outLength == 0 {
		outLength++
	}
	if !c.positive {
		// negative sign
		outLength++
	}
	hasFraction := c.fracSignificantTo > c.fracSignificantFrom
	if hasFraction {
		outLength += 1 + c.fracLeadingZeroCount + c.fracSignificantTo - c.fracSignificantFrom
	}

	/*
		// safety check for static buffer
		if len(buffer) < outLength {
			return "", false
		}
	*/

	c.builder.Reset()
	c.builder.Grow(outLength)

	if !c.positive {
		c.builder.WriteByte(minusByte)
	}

	if c.intSignificantFrom < c.intTo {
		// there is an integer part
		if c.positive {
			for i := c.intSignificantFrom; i < c.intSignificantTo; i++ {
				c.builder.WriteByte(c.input[i])
			}
		} else {
			for i := c.intSignificantFrom; i < c.intSignificantTo; i++ {
				c.builder.WriteByte(flipDigit(c.input[i]))
			}
		}

		// there are trailing zeros
		if c.intTo > c.intSignificantTo {
			for i := c.intSignificantTo; i < c.intTo; i++ {
				c.builder.WriteByte(digit0)
			}
		}
	} else {
		// there is no integer part
		c.builder.WriteByte(digit0)
	}

	if hasFraction {
		c.builder.WriteByte(decimalPoint)
		for i := 0; i < c.fracLeadingZeroCount; i++ {
			c.builder.WriteByte(digit0)
		}
		if c.positive {
			for i := c.fracSignificantFrom; i < c.fracSignificantTo; i++ {
				c.builder.WriteByte(c.input[i])
			}
		} else {
			for i := c.fracSignificantFrom; i < c.fracSignificantTo; i++ {
				c.builder.WriteByte(flipDigit(c.input[i]))
			}
		}
	}

	out = c.builder.String()
	ok = true
	c.input = ""

	// TODO remove debug
	if outLength != len(out) {
		panic(fmt.Sprintf("outLength error: for %v -> %v expected %d, got %d", s, out, outLength, len(out)))
	}

	return
}

// AnalyzeToken produces the correct internal state for the decoding step
func (c *codec) AnalyzeToken(s string) {
	c.cursor = 0
	c.length = len(c.input)
	c.ok = false

	if c.length == 0 {
		c.empty = true
		c.ok = true
		return
	}
	c.empty = false

	if c.input == zeroOutput {
		c.zero = true
		c.ok = true
		return
	}
	c.zero = false

	c.decodeSign(c.input[0])
	if !c.positive {
		if c.input[c.length-1] != negativeNumberTerminator {
			// negative terminator is not at the end as it should be
			return
		}
		// ignore the negative terminator from here on
		c.length--
	}

	if c.length < 3 {
		// too short to be valid
		return
	}

	c.magnitude = 0
	for c.cursor = 1; c.cursorCanRead(); c.cursor++ {
		m := c.decodeMagnitude(c.input[c.cursor])
		if m > maxMagnitudeDigitValue {
			c.magnitude += maxMagnitudeDigitValue
		} else {
			c.magnitude += m
			break
		}
	}
	c.cursor++

	c.ok = true

	if c.magnitudePositive {
		c.resetFractionParams()
		c.intSignificantFrom = c.cursor
		c.intTo = c.intSignificantFrom + c.magnitude
		if c.length <= c.intTo {
			// is an integer
			c.intSignificantTo = c.length
			return
		}
		// has fraction part too
		c.intSignificantTo = c.intTo
		c.fracSignificantFrom = c.intTo
		c.fracSignificantTo = c.length
		return
	}
	// is purely fractional
	c.intSignificantFrom = 0
	c.intSignificantTo = 0
	c.intTo = 0
	c.fracLeadingZeroCount = c.magnitude
	c.fracSignificantFrom = c.cursor
	c.fracSignificantTo = c.length
}

func (c *codec) EncodeInText(in string) (out string, ok bool) {
	inNum := false
	donePart := 0
	var b strings.Builder
	ok = true
	b.Grow(len(in))

	for i := 0; i < len(in); i++ {
		if in[i] >= digit0 && in[i] <= digit9 {
			if !inNum {
				b.Write([]byte(in[donePart:i]))
				donePart = i
				inNum = true
			}
			continue
		}
		if inNum {
			encoded, encOk := c.Encode(in[donePart:i])
			if encOk {
				b.Write([]byte(encoded))
			} else {
				b.WriteByte(negativeNumberTerminator)
				b.Write([]byte(in[donePart:i]))
				b.WriteByte(negativeNumberTerminator)
				ok = false
			}
			inNum = false
			donePart = i
		}
	}
	if !inNum {
		b.Write([]byte(in[donePart:]))
	} else {
		encoded, encOk := c.Encode(in[donePart:])
		if encOk {
			b.Write([]byte(encoded))
		} else {
			b.WriteByte(negativeNumberTerminator)
			b.Write([]byte(in[donePart:]))
			b.WriteByte(negativeNumberTerminator)
			ok = false
		}
	}

	out = b.String()
	return
}

func (c *codec) cursorCanRead() bool {
	return c.cursor < c.length
}

func (c *codec) checkSign() (found bool) {
	switch c.input[0] {
	case minusByte:
		c.positive = false
		found = true
	case plusByte:
		c.positive = true
		found = true
	default:
		c.positive = true
		found = false
	}
	return
}

func (c *codec) skipZeroes() {
	for ; c.cursorCanRead() && c.input[c.cursor] == digit0; c.cursor++ {
	}
}

func (c *codec) getIntPartBounds() (ok bool) {
	c.intSignificantFrom = c.cursor
	trailingZeroCount := 0
	for ; c.cursorCanRead(); c.cursor++ {
		if c.input[c.cursor] == decimalPoint {
			break
		}
		if c.input[c.cursor] == digit0 {
			trailingZeroCount++
			continue
		}
		if !isDigit(c.input[c.cursor]) {
			// some unexpected character encountered
			return false
		}
		if trailingZeroCount != 0 {
			trailingZeroCount = 0
		}
	}
	c.intSignificantTo = c.cursor - trailingZeroCount
	c.intTo = c.cursor

	if c.intSignificantFrom < c.intTo {
		c.zero = false
	}
	return true
}

func (c *codec) resetFractionParams() {
	c.fracSignificantFrom = 0
	c.fracSignificantTo = 0
	c.fracLeadingZeroCount = 0
}

func (c *codec) getFractionPartBounds() (ok bool) {
	fractionFrom := c.cursor
	c.skipZeroes()

	// fraction contains only zeroes and thus is ignored
	if !c.cursorCanRead() {
		return true
	}
	c.zero = false

	c.fracLeadingZeroCount = c.cursor - fractionFrom
	c.fracSignificantFrom = c.cursor
	trailingZeroCount := 0
	for ; c.cursorCanRead(); c.cursor++ {
		if c.input[c.cursor] == digit0 {
			trailingZeroCount++
			continue
		}
		if !isDigit(c.input[c.cursor]) {
			// some bogus character encountered
			return false
		}
		if trailingZeroCount != 0 {
			trailingZeroCount = 0
		}
	}
	c.fracSignificantTo = c.cursor - trailingZeroCount
	return true
}

func (c *codec) encodeSign() byte {
	if c.positive {
		if c.magnitudePositive {
			return signPositiveMagPositive
		}
		return signPositiveMagNegative
	}
	if c.magnitudePositive {
		return signNegativeMagPositive
	}
	return signNegativeMagNegative
}

func (c *codec) decodeSign(sign byte) {
	switch sign {
	case signPositiveMagPositive:
		c.positive = true
		c.magnitudePositive = true
	case signPositiveMagNegative:
		c.positive = true
		c.magnitudePositive = false
	case signNegativeMagNegative:
		c.positive = false
		c.magnitudePositive = false
	case signNegativeMagPositive:
		c.positive = false
		c.magnitudePositive = true
	}
}

func (c *codec) encodeMagnitude(m int) byte {
	if c.positive == c.magnitudePositive {
		return intToDigit(m)
	}
	return intToReversedDigit(m)
}

func (c *codec) decodeMagnitude(d byte) int {
	if c.positive == c.magnitudePositive {
		return digitToInt(d)
	}
	return reversedDigitToInt(d)
}
