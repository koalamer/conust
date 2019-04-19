package conust

import (
	"fmt"
	"strings"
)

// codec is the base 10 variant of the encoder
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
	intNonZeroFrom       int
	intNonZeroTo         int
	intTo                int
	fracLeadingZeroCount int
	fracNonZeroFrom      int
	fracNonZeroTo        int

	// output builder
	builder strings.Builder
}

// NewCodec returns a codec
func NewCodec() Codec {
	return &codec{}
}

//Encode turns input number into a base(10) Conust string
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

	intNonZeroLength := c.intNonZeroTo - c.intNonZeroFrom

	// determine allocation size
	outLength := 2 + intNonZeroLength // sign, magnitude, int length
	if intNonZeroLength == 0 {        // if int part is empty, that will be a single character of 0
		outLength++
	}
	if c.fracNonZeroTo-c.fracNonZeroFrom > 0 {
		outLength += 2 + (c.fracNonZeroTo - c.fracNonZeroFrom) // fraction plus separator and leading zero count character
	}
	if !c.positive { // will have single char postfix
		outLength++
	}

	// build encoded string
	c.builder.Reset()
	c.builder.Grow(outLength)

	if c.positive {
		c.builder.WriteByte(signPositive)
	} else {
		c.builder.WriteByte(signNegative)
	}

	if intNonZeroLength == 0 {
		if c.positive {
			c.builder.WriteByte(intToDigit(1))
			c.builder.WriteByte(digit0)
		} else {
			c.builder.WriteByte(intToReversedDigit(1))
			c.builder.WriteByte(intToReversedDigit(0))
		}
	} else {
		if c.positive {
			c.builder.WriteByte(intToDigit(c.intTo - c.intNonZeroFrom))
			for i := c.intNonZeroFrom; i < c.intNonZeroTo; i++ {
				c.builder.WriteByte(s[i])
			}
		} else {
			c.builder.WriteByte(intToReversedDigit(c.intTo - c.intNonZeroFrom))
			for i := c.intNonZeroFrom; i < c.intNonZeroTo; i++ {
				c.builder.WriteByte(flipDigit(s[i]))
			}
		}
	}

	if c.fracNonZeroTo > c.fracNonZeroFrom {
		if c.positive {
			c.builder.WriteByte(positiveIntegerTerminator)
			c.builder.WriteByte(intToReversedDigit(c.fracLeadingZeroCount))
			for i := c.fracNonZeroFrom; i < c.fracNonZeroTo; i++ {
				c.builder.WriteByte(s[i])
			}
		} else {
			c.builder.WriteByte(negativeIntegerTerminator)
			c.builder.WriteByte(intToDigit(c.fracLeadingZeroCount))
			for i := c.fracNonZeroFrom; i < c.fracNonZeroTo; i++ {
				c.builder.WriteByte(flipDigit(s[i]))
			}
			c.builder.WriteByte(negativeIntegerTerminator)
		}
	} else {
		if !c.positive {
			c.builder.WriteByte(negativeIntegerTerminator)
		}
	}
	out = c.builder.String()
	ok = true
	c.input = ""
	return
}

// Decode turns a Conust string back into its sourceBase representation
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

	outLength := c.intTo - c.intNonZeroFrom
	if !c.positive {
		outLength++
	}
	if c.fracNonZeroTo > c.fracNonZeroFrom {
		outLength += c.fracLeadingZeroCount
	}
	hasFraction := c.fracNonZeroTo > c.fracNonZeroFrom
	if hasFraction {
		outLength += 1 + c.fracNonZeroTo - c.fracNonZeroFrom
	}

	c.builder.Reset()
	c.builder.Grow(outLength)

	if c.positive {
		for i := c.intNonZeroFrom; i < c.intNonZeroTo; i++ {
			c.builder.WriteByte(c.input[i])
		}
	} else {
		c.builder.WriteByte(minusByte)
		for i := c.intNonZeroFrom; i < c.intNonZeroTo; i++ {
			c.builder.WriteByte(flipDigit(c.input[i]))
		}
	}

	if c.intTo > c.intNonZeroTo {
		for i := c.intNonZeroTo; i < c.intTo; i++ {
			c.builder.WriteByte(digit0)
		}
	}

	if hasFraction {
		if c.positive {
			c.builder.WriteByte(positiveIntegerTerminator)
		} else {
			c.builder.WriteByte(negativeIntegerTerminator)
		}
		for i := 0; i < c.fracLeadingZeroCount; i++ {
			c.builder.WriteByte(digit0)
		}
		if c.positive {
			for i := c.fracNonZeroFrom; i < c.fracNonZeroTo; i++ {
				c.builder.WriteByte(c.input[i])
			}
		} else {
			for i := c.fracNonZeroFrom; i < c.fracNonZeroTo; i++ {
				c.builder.WriteByte(flipDigit(c.input[i]))
			}
		}
	}

	out = c.builder.String()
	ok = true
	// TODO: remove debug
	if len(out) != outLength {
		fmt.Printf("out length miscalculation: encoded: %s, text: %s, length: %d, estimated %d\n", s, out, len(out), outLength)
	}
	c.input = ""
	return
}

func (c *codec) AnalyzeToken(s string) {
	/* defer func() {
		fmt.Printf("Token: %q, ok: %v, empty: %-5v, zero: %-5v, positive: %-5v\n", s, ok, c.isEmpty, c.isZero, c.isPositive)
		fmt.Printf("  int nz start: %-3d, nz end: %-3d, end: %-3d\n", c.intNonZeroFrom, c.intNonZeroTo, c.intTo)
		fmt.Printf("  frac leading z count: %-3d, nz start: %-3d, nz end: %-3d\n", c.fracLeadingZeroCount, c.fracNonZeroFrom, c.fracNonZeroTo)
		fmt.Printf("  base %d -> %d, precision %d -> %d, conversion needed: %v\n", c.sourceBase, c.newBase, c.sourceBasePrecision, c.newBasePrecision, c.baseConversionNeeded)
	}() */
	c.cursor = 0
	c.length = len(c.input)

	c.ok = true
	if !c.cursorValid() {
		c.empty = true
		return
	}
	c.empty = false

	if c.input == zeroOutput {
		c.zero = true
		return
	}
	c.zero = false

	c.checkTokenSign()
	var intTerminator byte
	var intLength int
	if c.positive {
		intTerminator = positiveIntegerTerminator
		intLength = digitToInt(s[1])
	} else {
		intTerminator = negativeIntegerTerminator
		intLength = reversedDigitToInt(s[1])
	}
	c.intNonZeroFrom = 2
	c.intTo = c.intNonZeroFrom + intLength

	terminatorPos := strings.IndexByte(s, intTerminator)
	c.initFractionParams()
	if terminatorPos < 0 {
		c.intNonZeroTo = c.length
		return
	}
	if terminatorPos <= c.intNonZeroFrom {
		c.ok = false
		return
	}
	c.intNonZeroTo = terminatorPos

	if c.positive {
		c.fracLeadingZeroCount = reversedDigitToInt(c.input[terminatorPos+1])
	} else {
		c.fracLeadingZeroCount = digitToInt(c.input[terminatorPos+1])
	}
	c.fracNonZeroFrom = terminatorPos + 2
	c.fracNonZeroTo = c.length
}

// AnalyzeInput disects identifies useful parts of the input and stores markers internally
func (c *codec) AnalyzeInput() {
	/* defer func() {
		fmt.Printf("Input: %q, ok: %v, empty: %-5v, zero: %-5v, positive: %-5v\n", s, ok, c.isEmpty, c.isZero, c.isPositive)
		fmt.Printf("  int nz start: %-3d, nz end: %-3d, end: %-3d\n", c.intNonZeroFrom, c.intNonZeroTo, c.intTo)
		fmt.Printf("  frac leading z count: %-3d, nz start: %-3d, nz end: %-3d\n", c.fracLeadingZeroCount, c.fracNonZeroFrom, c.fracNonZeroTo)
		fmt.Printf("  base %d -> %d, precision %d -> %d, conversion needed: %v\n", c.sourceBase, c.newBase, c.sourceBasePrecision, c.newBasePrecision, c.baseConversionNeeded)
	}() */
	c.cursor = 0
	c.length = len(c.input)

	// empty input results in empty but ok output
	if !c.cursorValid() {
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
	if !c.cursorValid() {
		c.ok = false
		return
	}

	// skip leading zeroes
	c.skipZeroes()

	// if there were only zeroes, the result is zeroOutput (c.zero = true)
	if !c.cursorValid() {
		c.ok = true
		return
	}

	// determine integer part bounds
	if !c.getIntPartBounds() {
		c.ok = false
		return
	}

	c.initFractionParams()

	// if no fraction present, end processing
	if !c.cursorValid() {
		c.ok = true
		return
	}
	// if the last non digit chaccter is not the decimal separator, that's an error
	if c.input[c.cursor] != positiveIntegerTerminator {
		c.ok = false
		return
	}
	// skip over decimal separator
	c.cursor++

	// process fraction part
	c.ok = c.getFractionPartBounds()
	return
}

func (c codec) cursorValid() bool {
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
	for ; c.cursorValid() && c.input[c.cursor] == digit0; c.cursor++ {
	}
}

func (c *codec) getIntPartBounds() (ok bool) {
	c.intNonZeroFrom = c.cursor
	trailingZeroCount := 0
	for ; c.cursorValid(); c.cursor++ {
		if c.input[c.cursor] == positiveIntegerTerminator {
			break
		}
		if c.input[c.cursor] == digit0 {
			trailingZeroCount++
			continue
		}
		if c.input[c.cursor] > digitZ || c.input[c.cursor] < digit0 ||
			(c.input[c.cursor] > digit9 && c.input[c.cursor] < digitA) {
			return false
		}
		if trailingZeroCount != 0 {
			trailingZeroCount = 0
		}
	}
	c.intNonZeroTo = c.cursor - trailingZeroCount
	c.intTo = c.cursor

	if c.intNonZeroFrom < c.intTo {
		c.zero = false
	}
	return true
}

func (c *codec) initFractionParams() {
	c.fracNonZeroFrom = 0
	c.fracNonZeroTo = 0
	c.fracLeadingZeroCount = 0
}

func (c *codec) getFractionPartBounds() (ok bool) {
	fractionFrom := c.cursor
	c.skipZeroes()

	// fraction contains only zeroes and thus is ignored
	if !c.cursorValid() {
		return true
	}
	c.zero = false

	c.fracLeadingZeroCount = c.cursor - fractionFrom
	c.fracNonZeroFrom = c.cursor
	trailingZeroCount := 0
	for ; c.cursorValid(); c.cursor++ {
		if c.input[c.cursor] == digit0 {
			trailingZeroCount++
			continue
		}
		if c.input[c.cursor] > digitZ || c.input[c.cursor] < digit0 ||
			(c.input[c.cursor] > digit9 && c.input[c.cursor] < digitA) {
			return false
		}
		if trailingZeroCount != 0 {
			trailingZeroCount = 0
		}
	}
	c.fracNonZeroTo = c.cursor - trailingZeroCount
	return true
}

func (c *codec) checkTokenSign() {
	switch c.input[0] {
	case signNegative:
		c.positive = false
	default:
		c.positive = true
	}
}
