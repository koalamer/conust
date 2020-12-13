package conust

import (
	"strings"
)

// Codec can transform strings to and from the Conust format.
//
// It has EncodeToken and DecodeToken functions to transform simple numbers to and from the Conust format.
//
// There is also EncodeMixedText, a convenience function, that encodes each group of decimal numbers
// and returns the resulting string. So that for example the strings "Item 20" and "Item 100" become
// "Item 722" and "Item 731" which sort as the numeric value in them would naturally imply.
type Codec struct {
	builder strings.Builder
}

// EncodeToken turns the input number into the alphanumerically sortable Conust string.
// If the input hase a base higher than 10 and contains letter characters, it must be lowercased.
// Note that if you want to incorporate the generated token into a string, and the token is not at
// the very end of it, then you will need to add a space character after the token to ensure correct
// sorting of the string.
// EncodeMixedText does that automatically
func (c *Codec) EncodeToken(input string) (out string, ok bool) {
	if input == "" {
		return "", true
	}

	if !c.isValidInput(input) {
		return "", false
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

// DecodeToken turns a Conust string back into its normal representation. The output will not reconstruct
// leading and trailing zeros. The plus sign for positive numbers is omitted as well.
func (c *Codec) DecodeToken(input string) (out string, ok bool) {
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
		if input[encodedLength-1] != negativeNumberTerminator {
			return "", false
		}

		encodedLength--
	}

	significantPartLength := encodedLength - sStartPos

	for i := sStartPos; i < encodedLength; i++ {
		if !isDigit(input[i]) {
			return "", false
		}
	}

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

// EncodeMixedText is a convinience function that replaces all groups of decimal numbers of the input
// with Conust strings also surrounding them with spaces (if not already present) to ensure the expected ordering
func (c *Codec) EncodeMixedText(input string) (out string, ok bool) {
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
			encoded, encOk := c.EncodeToken(input[donePartEnd:i])
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
		encoded, encOk := c.EncodeToken(input[donePartEnd:])
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

func (c *Codec) isValidInput(input string) bool {
	if !isSignByte(input[0]) && !isDigit(input[0]) {
		return false
	}

	decimalPointAlreadyFound := false
	for i := 1; i < len(input); i++ {
		if isDigit(input[i]) {
			continue
		}

		if decimalPointAlreadyFound {
			return false
		}

		if input[i] == decimalPoint {
			decimalPointAlreadyFound = true
			continue
		}

		return false
	}

	return true
}

func (c *Codec) getPositivity(input string) (positive bool) {
	return input[0] != minusByte
}

func (c *Codec) getSignificantStartPos(input string) int {
	i := 0
	for ; i < len(input); i++ {
		if isDigit(input[i]) && input[i] != digit0 {
			return i
		}
	}
	return -1
}

func (c *Codec) getSignificantEndPos(input string) int {
	i := len(input) - 1
	for ; i >= 0; i-- {
		if isDigit(input[i]) && input[i] != digit0 {
			return i + 1
		}
	}
	return -1
}

func (c *Codec) getDecimalPointPos(input string) int {
	return strings.IndexByte(input, decimalPoint)
}

func (c *Codec) getMagnitudeParams(inputLength int, sStartPos int, sEndPos int, decimalPointPos int) (magnitude int, magnitudePositive bool) {
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

func (c *Codec) calculateEncodedSize(positive bool, magnitude int, sStartPos int, sEndPos int, decimalPointPos int) int {
	length := 2 + (magnitude / maxMagnitudeDigitValue) + sEndPos - sStartPos
	if !positive {
		length++
	}
	if sStartPos < decimalPointPos && decimalPointPos < sEndPos {
		length--
	}
	return length
}

func (c *Codec) encodeSign(positive bool, magnitudePositive bool) byte {
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

func (c *Codec) writeMagnitude(positive bool, magnitudePositive bool, magnitude int) {
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

func (c *Codec) writeDigits(positive bool, digits string) {
	if positive {
		c.builder.WriteString(digits)
	} else {
		for i := 0; i < len(digits); i++ {
			c.builder.WriteByte(reverseDigit(digits[i]))
		}
	}
}

func (c *Codec) decodeSigns(in string) (positive bool, magnitudePositive bool, ok bool) {
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

func (c *Codec) decodeMagnitude(in string, positive bool, magnitudePositive bool) (magnitude int, significantPartPos int, ok bool) {
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

func (c *Codec) calculateDecodedLength(positive bool, magnitudePositive bool, magnitude int, significantPartLength int) int {
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
