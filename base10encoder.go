package conust

import (
	"fmt"
	"strings"
)

// Base10Encoder is the base 10 variant of the encoder
type Base10Encoder struct {
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
	e.AnalyzeInput(s)
	if !e.ok {
		return "", false
	}
	if e.empty {
		return "", true
	}
	if e.zero {
		return zeroOutput, true
	}

	intNonZeroLength := e.intNonZeroTo - e.intNonZeroFrom

	// determine allocation size
	outLength := 2 + (e.intNonZeroTo - e.intNonZeroFrom) // sign, magnitude, int length
	if intNonZeroLength == 0 {                           // if int part is empty, that will be a single character of 0
		outLength++
	}
	if e.fracNonZeroTo-e.fracNonZeroFrom > 0 {
		outLength += (e.fracNonZeroTo - e.fracNonZeroFrom) + 2 // fraction plus separator and leading zero count character
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
			e.builder.WriteByte(flipDigit10(intToDigit(1)))
			e.builder.WriteByte(flipDigit10(digit0))
		}
	} else {
		if e.positive {
			e.builder.WriteByte(intToDigit(e.intEnd - e.intNonZeroFrom + 1))
			for i := e.intNonZeroFrom; i < e.intNonZeroTo; i++ {
				e.builder.WriteByte(s[i])
			}
		} else {
			e.builder.WriteByte(flipDigit10(intToDigit(e.intEnd - e.intNonZeroFrom + 1)))
			for i := e.intNonZeroFrom; i < e.intNonZeroTo; i++ {
				e.builder.WriteByte(flipDigit10(s[i]))
			}
		}
	}

	if e.fracNonZeroTo > e.fracNonZeroFrom {
		if e.positive {
			e.builder.WriteByte(positiveIntegerTerminator)
			e.builder.WriteByte(flipDigit10(intToDigit(e.fracLeadingZeroCount)))
			for i := e.fracNonZeroFrom; i < e.fracNonZeroTo; i++ {
				e.builder.WriteByte(s[i])
			}
		} else {
			e.builder.WriteByte(negativeIntegerTerminator)
			e.builder.WriteByte(intToDigit(e.fracLeadingZeroCount))
			for i := e.fracNonZeroFrom; i < e.fracNonZeroTo; i++ {
				e.builder.WriteByte(flipDigit10(s[i]))
			}
			e.builder.WriteByte(negativeIntegerTerminator)

		}
	}
	out = e.builder.String()
	ok = true
	// debug
	if len(out) != outLength {
		fmt.Printf("out length miscalculation: text: %s, encoded: %s, length: %d, estimated %d\n", s, out, len(out), outLength)
	}
	return
}

// AnalyzeInput disects identifies useful parts of the input and stores markers internally
func (e *Base10Encoder) AnalyzeInput(s string) {
	/* defer func() {
		fmt.Printf("Input: %q, ok: %v, empty: %-5v, zero: %-5v, positive: %-5v\n", s, ok, e.isEmpty, e.isZero, e.isPositive)
		fmt.Printf("  int nz start: %-3d, nz end: %-3d, end: %-3d\n", e.intNonZeroFrom, e.intNonZeroTo, e.intEnd)
		fmt.Printf("  frac leading z count: %-3d, nz start: %-3d, nz end: %-3d\n", e.fracLeadingZeroCount, e.fracNonZeroFrom, e.fracNonZeroTo)
	}() */
	length := len(s)
	// empty input results in empty but ok output
	if length == 0 {
		e.ok = true
		e.empty = true
		return
	}

	e.empty = false
	e.zero = true
	i := 0

	// determine sign
	switch s[0] {
	case minusByte:
		e.positive = false
		i++
	case plusByte:
		e.positive = true
		i++
	default:
		e.positive = true
	}
	// a sign only is bad input
	if i >= length {
		e.ok = false
		return
	}

	// skip leading zeroes
	for ; i < length && s[i] == digit0; i++ {
	}

	// if there were only zeroes, the result is zeroOutput
	if i >= length {
		e.ok = true
		return
	}

	// determine integer part bounds
	e.intNonZeroFrom = i
	trailingZeroCount := 0
	for ; i < length; i++ {
		if s[i] == positiveIntegerTerminator {
			break
		}
		if s[i] == digit0 {
			trailingZeroCount++
			continue
		}
		if s[i] > digit9 || s[i] < digit0 {
			e.ok = false
			return
		}
		if trailingZeroCount != 0 {
			trailingZeroCount = 0
		}
	}
	e.intNonZeroTo = i - trailingZeroCount
	e.intEnd = (i - 1)
	if e.intNonZeroFrom <= e.intEnd {
		e.zero = false
	}

	// init fraction variables
	e.fracNonZeroFrom = 0
	e.fracNonZeroTo = 0
	e.fracLeadingZeroCount = 0

	// if no fraction present, end processing
	if i >= length-1 {
		e.ok = true
		return
	}

	// skip over decimal separator
	i++

	// process fraction part
	e.fracNonZeroFrom = i
	for ; i < length && s[i] == digit0; i++ {
	}

	// fraction contains only zeroes
	if i >= length {
		e.fracNonZeroFrom = 0
		e.ok = true
		return
	}

	e.zero = false
	e.fracLeadingZeroCount = i - e.fracNonZeroFrom
	e.fracNonZeroFrom = i
	trailingZeroCount = 0
	for ; i < length; i++ {
		if s[i] == digit0 {
			trailingZeroCount++
			continue
		}
		if s[i] > digit9 || s[i] < digit0 {
			e.ok = false
			return
		}
		if trailingZeroCount != 0 {
			trailingZeroCount = 0
		}
	}
	e.fracNonZeroTo = i - trailingZeroCount
	e.ok = true
	return
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
			b.WriteByte(flipDigit10(number[j]))
		}
		b.WriteByte(negativeIntegerTerminator)
	}
}
