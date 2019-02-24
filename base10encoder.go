package conust

import (
	"fmt"
	"strings"
)

// Base10Encoder is the base 10 variant of the encoder
type Base10Encoder struct {
	input                string
	length               int
	cursor               int
	ok                   bool
	empty                bool
	zero                 bool
	positive             bool
	intNonZeroFrom       int
	intNonZeroTo         int
	intEnd               int
	fracLeadingZeroCount int
	fracNonZeroFrom      int
	fracNonZeroTo        int
	builder              strings.Builder
}

// NewBase10Encoder returns an encoder that outputs base(10) Conust strings
func NewBase10Encoder() Encoder {
	return &Base10Encoder{}
}

//Encode turns input number into a base(10) Conust string
func (e *Base10Encoder) Encode(s string) (out string, ok bool) {
	e.input = s
	e.AnalyzeInput()
	if !e.ok {
		e.input = ""
		return "", false
	}
	if e.empty {
		e.input = ""
		return "", true
	}
	if e.zero {
		e.input = ""
		return zeroOutput, true
	}

	intNonZeroLength := e.intNonZeroTo - e.intNonZeroFrom

	// determine allocation size
	outLength := 2 + intNonZeroLength // sign, magnitude, int length
	if intNonZeroLength == 0 {        // if int part is empty, that will be a single character of 0
		outLength++
	}
	if e.fracNonZeroTo-e.fracNonZeroFrom > 0 {
		outLength += 2 + (e.fracNonZeroTo - e.fracNonZeroFrom) // fraction plus separator and leading zero count character
	}
	if !e.positive { // will have single char postfix
		outLength++
	}

	// build encoded string
	e.builder.Reset()
	e.builder.Grow(outLength)

	if e.positive {
		e.builder.WriteByte(signPositive)
	} else {
		e.builder.WriteByte(signNegative)
	}

	if intNonZeroLength == 0 {
		if e.positive {
			e.builder.WriteByte(intToDigit(1))
			e.builder.WriteByte(digit0)
		} else {
			e.builder.WriteByte(flipDigit36(intToDigit(1)))
			e.builder.WriteByte(flipDigit36(digit0))
		}
	} else {
		if e.positive {
			e.builder.WriteByte(intToDigit(e.intEnd - e.intNonZeroFrom + 1))
			for i := e.intNonZeroFrom; i < e.intNonZeroTo; i++ {
				e.builder.WriteByte(s[i])
			}
		} else {
			e.builder.WriteByte(flipDigit36(intToDigit(e.intEnd - e.intNonZeroFrom + 1)))
			for i := e.intNonZeroFrom; i < e.intNonZeroTo; i++ {
				e.builder.WriteByte(flipDigit36(s[i]))
			}
		}
	}

	if e.fracNonZeroTo > e.fracNonZeroFrom {
		if e.positive {
			e.builder.WriteByte(positiveIntegerTerminator)
			e.builder.WriteByte(flipDigit36(intToDigit(e.fracLeadingZeroCount)))
			for i := e.fracNonZeroFrom; i < e.fracNonZeroTo; i++ {
				e.builder.WriteByte(s[i])
			}
		} else {
			e.builder.WriteByte(negativeIntegerTerminator)
			e.builder.WriteByte(intToDigit(e.fracLeadingZeroCount))
			for i := e.fracNonZeroFrom; i < e.fracNonZeroTo; i++ {
				e.builder.WriteByte(flipDigit36(s[i]))
			}
			e.builder.WriteByte(negativeIntegerTerminator)
		}
	} else {
		if !e.positive {
			e.builder.WriteByte(negativeIntegerTerminator)
		}
	}
	out = e.builder.String()
	ok = true
	// debug
	if len(out) != outLength {
		fmt.Printf("out length miscalculation: text: %s, encoded: %s, length: %d, estimated %d\n", s, out, len(out), outLength)
	}
	e.input = ""
	return
}

// AnalyzeInput disects identifies useful parts of the input and stores markers internally
func (e *Base10Encoder) AnalyzeInput() {
	/* defer func() {
		fmt.Printf("Input: %q, ok: %v, empty: %-5v, zero: %-5v, positive: %-5v\n", s, ok, e.isEmpty, e.isZero, e.isPositive)
		fmt.Printf("  int nz start: %-3d, nz end: %-3d, end: %-3d\n", e.intNonZeroFrom, e.intNonZeroTo, e.intEnd)
		fmt.Printf("  frac leading z count: %-3d, nz start: %-3d, nz end: %-3d\n", e.fracLeadingZeroCount, e.fracNonZeroFrom, e.fracNonZeroTo)
	}() */
	e.cursor = 0
	e.length = len(e.input)

	// empty input results in empty but ok output
	if !e.cursorValid() {
		e.ok = true
		e.empty = true
		return
	}

	e.empty = false
	e.zero = true

	if e.checkSign() {
		e.cursor++
	}
	// a sign only is bad input
	if !e.cursorValid() {
		e.ok = false
		return
	}

	// skip leading zeroes
	e.skipZeroes()

	// if there were only zeroes, the result is zeroOutput (e.zero = true)
	if !e.cursorValid() {
		e.ok = true
		return
	}

	// determine integer part bounds
	if !e.getIntPartBounds() {
		e.ok = false
		return
	}

	e.initFractionParams()

	// if no fraction present, end processing
	if !e.cursorValid() {
		e.ok = true
		return
	}
	// if the last non digit chaccter is not the decimal separator, that's an error
	if e.input[e.cursor] != positiveIntegerTerminator {
		e.ok = false
		return
	}
	// skip over decimal separator
	e.cursor++

	// process fraction part
	e.ok = e.getFractionPartBounds()
	return
}

func (e Base10Encoder) cursorValid() bool {
	return e.cursor < e.length
}

func (e *Base10Encoder) checkSign() (found bool) {
	switch e.input[0] {
	case minusByte:
		e.positive = false
		found = true
	case plusByte:
		e.positive = true
		found = true
	default:
		e.positive = true
		found = false
	}
	return
}

func (e *Base10Encoder) skipZeroes() {
	for ; e.cursorValid() && e.input[e.cursor] == digit0; e.cursor++ {
	}
}

func (e *Base10Encoder) getIntPartBounds() (ok bool) {
	e.intNonZeroFrom = e.cursor
	trailingZeroCount := 0
	for ; e.cursorValid(); e.cursor++ {
		if e.input[e.cursor] == positiveIntegerTerminator {
			break
		}
		if e.input[e.cursor] == digit0 {
			trailingZeroCount++
			continue
		}
		if e.input[e.cursor] > digitZ || e.input[e.cursor] < digit0 ||
			(e.input[e.cursor] > digit9 && e.input[e.cursor] < digitA) {
			return false
		}
		if trailingZeroCount != 0 {
			trailingZeroCount = 0
		}
	}
	e.intNonZeroTo = e.cursor - trailingZeroCount
	e.intEnd = (e.cursor - 1)

	if e.intNonZeroFrom <= e.intEnd {
		e.zero = false
	}
	return true
}

func (e *Base10Encoder) initFractionParams() {
	e.fracNonZeroFrom = 0
	e.fracNonZeroTo = 0
	e.fracLeadingZeroCount = 0
}

func (e *Base10Encoder) getFractionPartBounds() (ok bool) {
	fractionFrom := e.cursor
	e.skipZeroes()

	// fraction contains only zeroes and thus is ignored
	if !e.cursorValid() {
		return true
	}
	e.zero = false

	e.fracLeadingZeroCount = e.cursor - fractionFrom
	e.fracNonZeroFrom = e.cursor
	trailingZeroCount := 0
	for ; e.cursorValid(); e.cursor++ {
		if e.input[e.cursor] == digit0 {
			trailingZeroCount++
			continue
		}
		if e.input[e.cursor] > digitZ || e.input[e.cursor] < digit0 ||
			(e.input[e.cursor] > digit9 && e.input[e.cursor] < digitA) {
			return false
		}
		if trailingZeroCount != 0 {
			trailingZeroCount = 0
		}
	}
	e.fracNonZeroTo = e.cursor - trailingZeroCount
	return true
}

// FromInt32 turns input number into a base(10) Conust string
func (e Base10Encoder) FromInt32(i int32) string {
	if i == 0 {
		return zeroOutput
	}
	return e.fromIntString(int32Preproc(i))
}

// FromInt64 turns input number into a base(10) Conust string
func (e Base10Encoder) FromInt64(i int64) string {
	if i == 0 {
		return zeroOutput
	}
	return e.fromIntString(int64Preproc(i))
}

// FromFloat32 turns input number into a base(10) Conust string
func (e Base10Encoder) FromFloat32(f float32) string {
	// TODO
	return ""
}

// FromFloat64 turns input number into a base(10) Conust string
func (e Base10Encoder) FromFloat64(f float64) string {
	// TODO
	return ""
}

func (e Base10Encoder) fromIntString(positive bool, absNumber string) string {
	var b strings.Builder
	b.Grow(builderInitialCap)
	if positive {
		b.WriteByte(signPositive)
		e.encode(&b, true, absNumber)
	} else {
		b.WriteByte(signNegative)
		e.encode(&b, false, absNumber)
	}
	return b.String()
}

func (e Base10Encoder) encode(b *strings.Builder, positive bool, number string) {
	if positive {
		b.WriteByte(intToDigit(len(number)))
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
