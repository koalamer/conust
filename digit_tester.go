package conust

import "errors"

type digitTester struct {
	minDecimalDigit   byte
	maxDecimalDigit   byte
	minLowercaseDigit byte
	maxLowercaseDigit byte
	minUppercaseDigit byte
	maxUppercaseDigit byte
}

const neverMatchingRangeMin = 1
const neverMatchingRangeMax = 0

var errIllegalRadix = errors.New("illegal radix")

func newDigitTester(radix int) (*digitTester, error) {
	if radix < 2 || radix > 36 {
		return nil, errIllegalRadix
	}

	if radix <= 10 {
		return &digitTester{
			minDecimalDigit:   digits36[0],
			maxDecimalDigit:   digits36[radix-1],
			minLowercaseDigit: neverMatchingRangeMin,
			maxLowercaseDigit: neverMatchingRangeMax,
			minUppercaseDigit: neverMatchingRangeMin,
			maxUppercaseDigit: neverMatchingRangeMax,
		}, nil
	}

	return &digitTester{
		minDecimalDigit:   digits36[0],
		maxDecimalDigit:   digits36[9],
		minLowercaseDigit: digits36[10],
		maxLowercaseDigit: digits36[radix-1],
		minUppercaseDigit: uppercaseDigits36[10],
		maxUppercaseDigit: uppercaseDigits36[radix-1],
	}, nil
}

func (dt *digitTester) isDigit(b byte) bool {
	return (b <= dt.maxDecimalDigit && b >= dt.minDecimalDigit) ||
		(b <= dt.maxLowercaseDigit &&
			(b >= dt.minLowercaseDigit ||
				(b <= dt.maxUppercaseDigit && b >= dt.minUppercaseDigit)))
}
