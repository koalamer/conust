package conust

import (
	"fmt"
	"strings"
)

// Base10Encoder is the base 10 variant of the encoder
type Base10Encoder struct {
	isEmpty              bool
	isZero               bool
	isPositive           bool
	intNonZeroFrom       int
	intNonZeroTo         int
	intEnd               int
	fracLeadingZeroCount int
	fracNonZeroFrom      int
	fracNonZeroTo        int
}

// NewBase10Encoder returns an encoder that outputs base(10) Conust strings
func NewBase10Encoder() Encoder {
	return &Base10Encoder{}
}

//FromString turns input number into a base(10) Conust string
func (e Base10Encoder) FromString(s string) (out string, ok bool) {
	// TODO
	return "", false
}

// AnalyzeInput disects identifies useful parts of the input and stores markers internally
func (e *Base10Encoder) AnalyzeInput(s string) (ok bool) {
	defer func() {
		fmt.Printf("Input: %q, ok: %v, empty: %-5v, zero: %-5v, positive: %-5v\n", s, ok, e.isEmpty, e.isZero, e.isPositive)
		fmt.Printf("  int nz start: %-3d, nz end: %-3d, end: %-3d\n", e.intNonZeroFrom, e.intNonZeroTo, e.intEnd)
		fmt.Printf("  frac leading z count: %-3d, nz start: %-3d, nz end: %-3d\n", e.fracLeadingZeroCount, e.fracNonZeroFrom, e.fracNonZeroTo)
	}()
	length := len(s)
	// empty input results in empty but ok output
	if length == 0 {
		e.isEmpty = true
		return true
	}

	e.isEmpty = false
	e.isZero = true
	i := 0

	// determine sign
	switch s[0] {
	case minusByte:
		e.isPositive = false
		i++
	case plusByte:
		e.isPositive = true
		i++
	default:
		e.isPositive = true
	}
	// a sign only is bad input
	if i >= length {
		return false
	}

	// skip leading zeroes
	for ; i < length && s[i] == digit0; i++ {
	}

	// if there were only zeroes, the result is zeroOutput
	if i >= length {
		return true
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
			return false
		}
		if trailingZeroCount != 0 {
			trailingZeroCount = 0
		}
	}
	e.intNonZeroTo = i - trailingZeroCount
	e.intEnd = (i - 1)
	if e.intNonZeroFrom <= e.intEnd {
		e.isZero = false
	}

	// init fraction variables
	e.fracNonZeroFrom = 0
	e.fracNonZeroTo = 0
	e.fracLeadingZeroCount = 0

	// if no fraction present, end processing
	if i >= length-1 {
		return true
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
		return true
	}

	e.isZero = false
	e.fracLeadingZeroCount = i - e.fracNonZeroFrom
	e.fracNonZeroFrom = i
	trailingZeroCount = 0
	for ; i < length; i++ {
		if s[i] == digit0 {
			trailingZeroCount++
			continue
		}
		if s[i] > digit9 || s[i] < digit0 {
			return false
		}
		if trailingZeroCount != 0 {
			trailingZeroCount = 0
		}
	}
	e.fracNonZeroTo = i - trailingZeroCount
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
			b.WriteByte(flipDigit10(number[j]))
		}
		b.WriteByte(negativeIntegerTerminator)
	}
}
