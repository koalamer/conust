package conust

import (
	"strings"
)

type codec struct {
	builder strings.Builder
}

// NewCodec creates a new codec instance. These instances are not thread safe.
func NewCodec() Codec {
	return &codec{}
}

func (c *codec) Encode(input string) (out string, ok bool) {
	return c.EncodeToken(input)
}

func (c *codec) EncodeToken(input string) (out string, ok bool) {
	if input == "" {
		return "", true
	}

	positive := c.getPositivity(input)
	decimalPointPos := c.getDecimalPointPos(input)
	sStartPos := c.getSignificantStartPos(input)
	sEndPos := c.getSignificantEndPos(input)

	if sStartPos == sEndPos {
		return zeroOutput, true
	}

	magnitude, magnitudePositive := c.getMagnitudeParams(len(input), sStartPos, sEndPos, decimalPointPos)

	c.builder.Reset()
	c.builder.Grow(c.calculateEncodedSize(positive, magnitude, sStartPos, sEndPos, decimalPointPos))
	c.builder.WriteByte(c.encodeSign(positive, magnitudePositive))
	c.writeMagnitude(positive, magnitudePositive, magnitude)

	if sStartPos < decimalPointPos && decimalPointPos < sEndPos {
		c.writeDigits(positive, input[sStartPos:decimalPointPos])
		c.writeDigits(positive, input[decimalPointPos+1:sEndPos])
	} else {
		c.writeDigits(positive, input[sStartPos:sEndPos])
	}
	if !positive {
		c.builder.WriteByte(negativeNumberTerminator)
	}
	return c.builder.String(), true
}

func (c *codec) Decode(input string) (out string, ok bool) {
	return c.DecodeToken(input)
}

func (c *codec) DecodeToken(input string) (out string, ok bool) {
	if input == "" {
		return "", true
	}

	if input == zeroOutput {
		return zeroInput, true
	}

	if len(input) < 3 {
		return "", false
	}

	positive, magnitudePositive, ok := c.decodeSigns(input)
	if !ok {
		return "", false
	}

	magnitude, sStartPos, ok := c.decodeMagnitude(input, positive, magnitudePositive)
	if !ok {
		return "", false
	}

	encodedLength := len(input)
	if !positive {
		encodedLength--
	}

	significantPartLength := encodedLength - sStartPos

	c.builder.Reset()
	c.builder.Grow(c.calculateDecodedLength(positive, magnitudePositive, magnitude, significantPartLength))

	if !positive {
		c.builder.WriteByte(minusByte)
	}
	if !magnitudePositive {
		c.builder.WriteByte(digit0)
		c.builder.WriteByte(decimalPoint)
		for i := 0; i < magnitude; i++ {
			c.builder.WriteByte(digit0)
		}
		c.writeDigits(positive, input[sStartPos:encodedLength])
	} else {
		if magnitude >= significantPartLength {
			c.writeDigits(positive, input[sStartPos:encodedLength])
			for i := 0; i < magnitude-significantPartLength; i++ {
				c.builder.WriteByte(digit0)
			}
		} else {
			c.writeDigits(positive, input[sStartPos:sStartPos+magnitude])
			c.builder.WriteByte(decimalPoint)
			c.writeDigits(positive, input[sStartPos+magnitude:encodedLength])
		}
	}

	return c.builder.String(), true
}

func (c *codec) EncodeMixedText(input string) (out string, ok bool) {
	insideNumber := false
	donePartEnd := 0
	var b strings.Builder
	ok = true
	b.Grow(len(input) + 6)

	for i := 0; i < len(input); i++ {
		if input[i] >= digit0 && input[i] <= digit9 {
			if !insideNumber {
				b.Write([]byte(input[donePartEnd:i]))
				donePartEnd = i
				insideNumber = true
				if i > 0 && input[i-1] != inTextSeparator {
					b.WriteByte(inTextSeparator)
				}
			}
			continue
		}
		if insideNumber {
			encoded, encOk := c.Encode(input[donePartEnd:i])
			if encOk {
				b.WriteString(encoded)
			} else {
				b.WriteString(input[donePartEnd:i])
				ok = false
			}
			insideNumber = false
			donePartEnd = i
			if input[i] != inTextSeparator {
				b.WriteByte(inTextSeparator)
			}
		}
	}
	if !insideNumber {
		b.WriteString(input[donePartEnd:])
	} else {
		encoded, encOk := c.Encode(input[donePartEnd:])
		if encOk {
			b.WriteString(encoded)
		} else {
			b.WriteString(input[donePartEnd:])
			ok = false
		}
	}

	out = b.String()
	return
}

func (c *codec) getPositivity(input string) (positive bool) {
	return input[0] != minusByte
}

func (c *codec) getSignificantStartPos(input string) int {
	i := 0
	for ; i < len(input); i++ {
		if isDigit(input[i]) && input[i] != digit0 {
			return i
		}
	}
	return -1
}

func (c *codec) getSignificantEndPos(input string) int {
	i := len(input) - 1
	for ; i >= 0; i-- {
		if isDigit(input[i]) && input[i] != digit0 {
			return i + 1
		}
	}
	return -1
}

func (c *codec) getDecimalPointPos(input string) int {
	return strings.IndexByte(input, decimalPoint)
}

func (c *codec) getMagnitudeParams(inputLength int, sStartPos int, sEndPos int, decimalPointPos int) (magnitude int, magnitudePositive bool) {
	if decimalPointPos < 0 {
		magnitude = inputLength - sStartPos
		magnitudePositive = true
	} else if decimalPointPos < sStartPos {
		magnitude = sStartPos - (decimalPointPos + 1)
		magnitudePositive = false
	} else {
		magnitude = decimalPointPos - sStartPos
		magnitudePositive = true
	}
	return
}

func (c *codec) calculateEncodedSize(positive bool, magnitude int, sStartPos int, sEndPos int, decimalPointPos int) int {
	length := 2 + (magnitude / maxMagnitudeDigitValue) + sEndPos - sStartPos
	if !positive {
		length++
	}
	if sStartPos < decimalPointPos && decimalPointPos < sEndPos {
		length--
	}
	return length
}

func (c *codec) encodeSign(positive bool, magnitudePositive bool) byte {
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

func (c *codec) writeMagnitude(positive bool, magnitudePositive bool, magnitude int) {
	reverseDigits := positive != magnitudePositive
	for ; magnitude > maxMagnitudeDigitValue; magnitude -= maxMagnitudeDigitValue {
		if reverseDigits {
			c.builder.WriteByte(intToReversedDigit(maxDigitValue))
		} else {
			c.builder.WriteByte(intToDigit(maxDigitValue))
		}
	}
	if reverseDigits {
		c.builder.WriteByte(intToReversedDigit(magnitude))
	} else {
		c.builder.WriteByte(intToDigit(magnitude))
	}
}

func (c *codec) writeDigits(positive bool, digits string) {
	if positive {
		c.builder.WriteString(digits)
	} else {
		for i := 0; i < len(digits); i++ {
			c.builder.WriteByte(reverseDigit(digits[i]))
		}

	}
}

func (c *codec) decodeSigns(in string) (positive bool, magnitudePositive bool, ok bool) {
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

func (c *codec) decodeMagnitude(in string, positive bool, magnitudePositive bool) (magnitude int, significantPartPos int, ok bool) {
	reverseDigits := positive != magnitudePositive
	var digitValue int
	for i := 1; i < len(in); i++ {
		if reverseDigits {
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

func (c *codec) calculateDecodedLength(positive bool, magnitudePositive bool, magnitude int, significantPartLength int) int {
	var signLength int
	if !positive {
		signLength = 1
	}
	if magnitudePositive {
		if magnitude >= significantPartLength {
			return signLength + magnitude
		}
		return signLength + significantPartLength + 1
	}
	return signLength + 2 + magnitude + significantPartLength
}
