package conust

import (
	"strconv"
	"strings"
)

type base36Encoder struct{}

// NewBase36Encoder returns a Decoder that outputs base(36) Conust strings
func NewBase36Encoder() Encoder {
	return base36Encoder{}
}

func (e base36Encoder) FromString(s string) (string, bool) {
	// TODO
	return "", false
}

func (e base36Encoder) AnalyzeInput(s string) {
	return
}

// FromInt32 turns input number into a base(36) Conust string
func (e base36Encoder) FromInt32(i int32) (s string) {
	return e.FromInt64(int64(i))
}

// FromInt64  turns input number into a base(36) Conust string
func (e base36Encoder) FromInt64(i int64) (s string) {
	if i == 0 {
		return zeroOutput
	}
	var b strings.Builder
	b.Grow(builderInitialCap)
	var number string
	if i > 0 {
		b.WriteByte(signPositive)
		number = strconv.FormatInt(i, 36)
		e.encode(&b, true, number)
	} else {
		b.WriteByte(signNegative)
		number = strconv.FormatInt(i, 36)[1:]
		e.encode(&b, false, number)
	}
	return b.String()
}

func (e base36Encoder) FromFloat32(f float32) string {
	// TODO
	return ""
}

func (e base36Encoder) FromFloat64(f float64) string {
	// TODO
	return ""
}
func (e base36Encoder) encode(b *strings.Builder, positive bool, number string) {
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
