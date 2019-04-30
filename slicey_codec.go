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

	positive := sc.getPositivity(in)
	decimalPointPos := sc.getDecimalPointPosition(in)
	start := sc.getStartPosition(in)
	end := sc.getEndPosition(in)

	if start == end {
		return zeroOutput, true
	}

	magnitude, magnitudePositive := sc.checkMagnitudeParams(len(in), start, end, decimalPointPos)

	sc.builder.Reset()
	sc.builder.Grow(sc.calculateEncodedSize(positive, magnitude, start, end, decimalPointPos))
	sc.builder.WriteByte(sc.encodeSign(positive, magnitudePositive))
	sc.writeMagnitude(positive, magnitudePositive, magnitude)

	if start < decimalPointPos && decimalPointPos < end {
		sc.writeDigits(positive, in[start:decimalPointPos])
		sc.writeDigits(positive, in[decimalPointPos+1:end])
	} else {
		sc.writeDigits(positive, in[start:end])
	}
	if !positive {
		sc.builder.WriteByte(negativeNumberTerminator)
	}
	return sc.builder.String(), true
}

func (sc *sliceyCodec) Decode(in string) (out string, ok bool) {
	if in == "" {
		return "", true
	}

	if in == zeroOutput {
		return zeroInput, true
	}

	if len(in) < 3 {
		return "", false
	}

	positive, magnitudePositive, ok := sc.getEncodedSigns(in)
	if !ok {
		return "", false
	}

	magnitude, significantPartPos, ok := sc.getEncodedMagnitude(in, positive, magnitudePositive)
	if !ok {
		return "", false
	}

	length := len(in)
	if !positive {
		length--
	}

	significantPartLength := length - significantPartPos
	decodedLength := sc.calculateDecodedLength(positive, magnitudePositive, magnitude, significantPartLength)

	sc.builder.Reset()
	sc.builder.Grow(decodedLength)

	if !positive {
		sc.builder.WriteByte(minusByte)
	}
	if !magnitudePositive {
		sc.builder.WriteByte(digit0)
		sc.builder.WriteByte(decimalPoint)
		for i := 0; i < magnitude; i++ {
			sc.builder.WriteByte(digit0)
		}
		sc.writeDigits(positive, in[significantPartPos:length])
	} else {
		if magnitude >= significantPartLength {
			sc.writeDigits(positive, in[significantPartPos:length])
			for i := 0; i < magnitude-significantPartLength; i++ {
				sc.builder.WriteByte(digit0)
			}
		} else {
			sc.writeDigits(positive, in[significantPartPos:significantPartPos+magnitude])
			sc.builder.WriteByte(decimalPoint)
			sc.writeDigits(positive, in[significantPartPos+magnitude:length])
		}
	}

	return sc.builder.String(), true
}

func (sc *sliceyCodec) EncodeInText(in string) (out string, ok bool) {
	inNum := false
	donePart := 0
	var b strings.Builder
	ok = true
	b.Grow(len(in) + 10)

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
			encoded, encOk := sc.Encode(in[donePart:i])
			if encOk {
				b.WriteString(encoded)
			} else {
				b.WriteString(in[donePart:i])
				ok = false
			}
			inNum = false
			donePart = i
		}
	}
	if !inNum {
		b.WriteString(in[donePart:])
	} else {
		encoded, encOk := sc.Encode(in[donePart:])
		if encOk {
			b.WriteString(encoded)
		} else {
			b.WriteString(in[donePart:])
			ok = false
		}
	}

	out = b.String()
	return
}

func (sc *sliceyCodec) getPositivity(in string) (positive bool) {
	return in[0] != minusByte
}

func (sc *sliceyCodec) getStartPosition(in string) (start int) {
	i := 0
	for ; i < len(in); i++ {
		if isDigit(in[i]) && in[i] != digit0 {
			return i
		}
	}
	return -1
}

func (sc *sliceyCodec) getEndPosition(in string) (end int) {
	i := len(in) - 1
	for ; i >= 0; i-- {
		if isDigit(in[i]) && in[i] != digit0 {
			return i + 1
		}
	}
	return -1
}

func (sc *sliceyCodec) getDecimalPointPosition(in string) (decimalPointPos int) {
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

func (sc *sliceyCodec) calculateEncodedSize(positive bool, magnitude int, start int, end int, decimalPointPos int) (encodedLength int) {
	encodedLength = 2 + (magnitude / maxMagnitudeDigitValue) + end - start
	if !positive {
		encodedLength++
	}
	if start < decimalPointPos && decimalPointPos < end {
		encodedLength--
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

func (sc *sliceyCodec) writeMagnitude(positive bool, magnitudePositive bool, magnitude int) {
	flippedDigits := positive != magnitudePositive
	for ; magnitude > maxMagnitudeDigitValue; magnitude -= maxMagnitudeDigitValue {
		if flippedDigits {
			sc.builder.WriteByte(intToReversedDigit(maxDigitValue))
		} else {
			sc.builder.WriteByte(intToDigit(maxDigitValue))
		}
	}
	if flippedDigits {
		sc.builder.WriteByte(intToReversedDigit(magnitude))
	} else {
		sc.builder.WriteByte(intToDigit(magnitude))
	}
}

func (sc *sliceyCodec) writeDigits(positive bool, digits string) {
	if positive {
		sc.builder.WriteString(digits)
	} else {
		for i := 0; i < len(digits); i++ {
			sc.builder.WriteByte(flipDigit(digits[i]))
		}

	}
}

func (sc *sliceyCodec) getEncodedSigns(in string) (positive bool, magnitudePositive bool, ok bool) {
	switch in[0] {
	case signPositiveMagPositive:
		return true, true, true
	case signPositiveMagNegative:
		return true, false, true
	case signNegativeMagNegative:
		return false, false, true
	case signNegativeMagPositive:
		return false, true, true
	default:
		return false, false, false
	}
}

// magnitude, significantPartPos, ok := sc.getEncodedMagnitude(in)
func (sc *sliceyCodec) getEncodedMagnitude(in string, positive bool, magnitudePositive bool) (magnitude int, significantPartPos int, ok bool) {
	flippedDigits := positive != magnitudePositive
	var digitValue int
	for i := 1; i < len(in); i++ {
		if flippedDigits {
			digitValue = reversedDigitToInt(in[i])
		} else {
			digitValue = digitToInt(in[i])
		}

		if digitValue == maxDigitValue {
			magnitude += maxMagnitudeDigitValue
		} else {
			magnitude += digitValue
			significantPartPos = i + 1
			ok = true
			return
		}
	}
	return 0, 0, false
}

func (sc *sliceyCodec) calculateDecodedLength(positive bool, magnitudePositive bool, magnitude int, significantPartLength int) (decodedLength int) {
	if !positive {
		decodedLength = 1
	}
	if magnitudePositive {
		if magnitude >= significantPartLength {
			decodedLength += magnitude
			return
		}
		decodedLength += significantPartLength + 1
		return
	}
	decodedLength += 2 + magnitude + significantPartLength
	return
}
