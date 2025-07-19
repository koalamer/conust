package conust

import "errors"

var errSameSeparator = errors.New("the two separators are the same")
var errInvalidDecimalSeparator = errors.New("the decimal separator cannot be 0-9, a-z, + or -")
var errInvalidThousandSeparator = errors.New("the thousand separator cannot be 0-9, a-z, + or -")

type Comparator struct {
	decimalSeparator  byte
	thousandSeparator byte
}

// New returns a *Comparator or an error if the parameters are wrong
func New(decimalSeparator, thousandSeparator byte) (comparator *Comparator, err error) {
	if decimalSeparator == thousandSeparator {
		return nil, errSameSeparator
	}

	if isSignByte(decimalSeparator) || isDigit(decimalSeparator) {
		return nil, errInvalidDecimalSeparator
	}

	if isSignByte(thousandSeparator) || isDigit(thousandSeparator) {
		return nil, errInvalidThousandSeparator
	}

	return &Comparator{
		decimalSeparator:  decimalSeparator,
		thousandSeparator: thousandSeparator,
	}, nil
}

// Compare returns -1 if a is less than b, 0 is a equals b, 1 if a is greater than b.
//
// If the numeric values equal, then string comparison will decide the result
// func (c *Comparator) Compare(a, b string) int {
// TODO
// }
