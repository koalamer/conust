package conust

import (
	"strings"
)

type sliceyCodec struct {
	builder strings.Builder
}

// NewSliceyCodec creates a slicey kind of codec
func NewSliceyCodec() (out Codec) {
	out = &sliceyCodec{}
	return
}

func (sc *sliceyCodec) Encode(in string) (out string, ok bool) {
	if in == "" {
		return "", true
	}

	positive := sc.checkPositive(in)
	decimalPointPos := sc.checkDecimalPointPos(in)
	start := sc.checkStart(in)
	end := sc.checkEnd(in)

	if start == end {
		return zeroOutput, true
	}

	magnitude, magnitudePositive := sc.checkMagnitudeParams(len(in), start, end, decimalPointPos)

	sc.builder.Reset()
	sc.builder.Grow(end - start + 5)
	sc.builder.WriteByte(sc.encodeSign(positive, magnitudePositive))
	sc.writeMagnitude(&sc.builder, positive, magnitudePositive, magnitude)

	if start < decimalPointPos && decimalPointPos < end {
		sc.writeDigits(&sc.builder, positive, in[start:decimalPointPos])
		sc.writeDigits(&sc.builder, positive, in[decimalPointPos+1:end])
	} else {
		sc.writeDigits(&sc.builder, positive, in[start:end])
	}
	if !positive {
		sc.builder.WriteByte(negativeNumberTerminator)
	}
	return sc.builder.String(), true
}

func (sc *sliceyCodec) Decode(in string) (out string, ok bool) {
	return "", false
}

func (sc *sliceyCodec) EncodeInText(in string) (out string, ok bool) {
	return "", false
}

func (sc *sliceyCodec) checkPositive(in string) (positive bool) {
	return in[0] != minusByte
}

func (sc *sliceyCodec) checkStart(in string) (start int) {
	found := false
	i := 0
	for ; i < len(in); i++ {
		if isDigit(in[i]) && in[i] != digit0 {
			found = true
			break
		}
	}
	if !found {
		return -1
	}
	return i
}

func (sc *sliceyCodec) checkEnd(in string) (end int) {
	found := false
	i := len(in) - 1
	for ; i >= 0; i-- {
		if isDigit(in[i]) && in[i] != digit0 {
			found = true
			break
		}
	}
	if !found {
		return -1
	}
	return i + 1
}

func (sc *sliceyCodec) checkDecimalPointPos(in string) (decimalPointPos int) {
	return strings.IndexByte(in, decimalPoint)
}

func (sc *sliceyCodec) checkMagnitudeParams(length int, start int, end int, decimalPointPos int) (magnitude int, magnitudePositive bool) {
	if decimalPointPos < 0 {
		magnitude = length - start
		magnitudePositive = true
	} else if decimalPointPos < start {
		magnitude = start - (decimalPointPos + 1)
		magnitudePositive = false
	} else {
		magnitude = decimalPointPos - start
		magnitudePositive = true
	}
	return
}

func (sc *sliceyCodec) encodeSign(positive bool, magnitudePositive bool) byte {
	if positive {
		if magnitudePositive {
			return signPositiveMagPositive
		}
		return signPositiveMagNegative
	}
	if magnitudePositive {
		return signNegativeMagPositive
	}
	return signNegativeMagNegative
}

func (sc *sliceyCodec) writeMagnitude(b *strings.Builder, positive bool, magnitudePositive bool, magnitude int) {
	for ; magnitude > maxMagnitudeDigitValue; magnitude -= maxMagnitudeDigitValue {
		b.WriteByte(sc.encodeMagnitude(maxDigitValue, positive, magnitudePositive))
	}
	b.WriteByte(sc.encodeMagnitude(magnitude, positive, magnitudePositive))
}

func (sc *sliceyCodec) encodeMagnitude(m int, positive bool, magnitudePositive bool) byte {
	if positive == magnitudePositive {
		return intToDigit(m)
	}
	return intToReversedDigit(m)
}

func (sc *sliceyCodec) writeDigits(b *strings.Builder, positive bool, digits string) {
	if positive {
		for i := 0; i < len(digits); i++ {
			b.WriteByte(digits[i])
		}
	} else {
		for i := 0; i < len(digits); i++ {
			b.WriteByte(flipDigit(digits[i]))
		}

	}

}
