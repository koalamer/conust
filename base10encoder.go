package conust

import (
	"strings"
)

type base10Encoder struct{}

// NewBase10Encoder returns an encoder that outputs base(10) Conust strings
func NewBase10Encoder() Encoder {
	return base10Encoder{}
}

//FromString turns input number into a base(10) Conust string
func (e base10Encoder) FromString(s string) string {
	//TODO
	return ""
}

// FromI32 turns input number into a base(10) Conust string
func (e base10Encoder) FromInt32(i int32) string {
	if i == 0 {
		return zeroOutput
	}
	return e.fromIntString(int32Preproc(i))
}

// FromInt64 turns input number into a base(10) Conust string
func (e base10Encoder) FromInt64(i int64) string {
	if i == 0 {
		return zeroOutput
	}
	return e.fromIntString(int64Preproc(i))
}

// FromFloat32 turns input number into a base(10) Conust string
func (e base10Encoder) FromFloat32(f float32) string {
	// TODO
	return ""
}

// FromFloat64 turns input number into a base(10) Conust string
func (e base10Encoder) FromFloat64(f float64) string {
	// TODO
	return ""
}

func (e base10Encoder) fromIntString(positive bool, absNumber string) string {
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

func (e base10Encoder) encode(b *strings.Builder, positive bool, number string) {
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
