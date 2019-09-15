package conust

import (
	"testing"
)

func TestArrayReversion(t *testing.T) {
	digitDictionaryLength := len(digits36)
	expectedDictionaryLength := 36

	if digitDictionaryLength != expectedDictionaryLength {
		t.Fatalf("Digit dictionary length is %d instead of %d", digitDictionaryLength, expectedDictionaryLength)
	}

	if digitDictionaryLength != len(digits36Reversed) {
		t.Fatal("Forward and backward digit dictionaries are of different length")
	}

	offset := digitDictionaryLength - 1
	for i, c := range digits36Reversed {
		if c != digits36[offset-i] {
			t.Fatalf("digit backward[%d] = %s but forward[%d] = %s",
				i,
				string(digits36Reversed[i]),
				offset-i,
				string(digits36[offset-i]),
			)
		}
	}
}

func TestArraySortedness(t *testing.T) {
	tests := []struct {
		name      string
		arr       []byte
		ascending bool
	}{
		{"digits36", digits36[:], true},
		{"digits36Reversed", digits36Reversed[:], false},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			for i := 1; i < len(testCase.arr); i++ {
				if (testCase.ascending && string(testCase.arr[i-1]) >= string(testCase.arr[i])) ||
					(!testCase.ascending && string(testCase.arr[i-1]) <= string(testCase.arr[i])) {
					t.Fatalf("%s has a sorting error between inexes %d and %d", testCase.name, i-1, i)
				}
			}
		})
	}
}

func TestSignBytes(t *testing.T) {
	if string(signNegativeMagPositive) >= string(signNegativeMagNegative) {
		t.Fatal("signNegativeMagPositive is not smaller than signNegativeMagNegative")
	}
	if string(signNegativeMagNegative) >= zeroOutput {
		t.Fatal("signNegativeMagNegative is not smaller than zeroOutput")
	}
	if zeroOutput >= string(signPositiveMagNegative) {
		t.Fatal("zeroOutput is not smaller than signPositiveMagNegative")
	}
	if string(signPositiveMagNegative) >= string(signPositiveMagPositive) {
		t.Fatal("signPositiveMagNegative is not smaller than signPositiveMagPositive")
	}
}

func TestNamedBytes(t *testing.T) {
	if digit0 != digits36[0] {
		t.Fatal("digit0 is not the [0] digit")
	}
	if digit9 != digits36[9] {
		t.Fatal("digit9 is not the [9] digit")
	}
	if digitA != digits36[10] {
		t.Fatal("digitA is not the [10] digit")
	}
	if digitZ != digits36[35] {
		t.Fatal("digitZ is not the [35] digit")
	}
}

func TestTerminatorBytes(t *testing.T) {
	if negativeNumberTerminator <= digits36[35] {
		t.Fatal("the negative decimal separator is not greater than the digits")
	}
	if inTextSeparator >= digits36[0] {
		t.Fatal("the in text separator is not smaller than the digits")
	}
}

func TestBoundaryVariables(t *testing.T) {
	if LessThanAny >= string(signNegativeMagPositive) {
		t.Fatal("the LessThanAny string is not smaller than the negative sign marker")
	}
	if GreaterThanAny <= string(signPositiveMagPositive) {
		t.Fatal("the GreaterThanAny string is not greater than the positive sign marker")
	}
}

func TestDigitValueLimits(t *testing.T) {
	if len(digits36)-1 != maxDigitValue {
		t.Fatal("maxDigitValue is not in sync with the digit dictionary size")
	}
	if maxDigitValue-1 != maxMagnitudeDigitValue {
		t.Fatal("maxMagnitudeDigitValue is not in sync with maxDigitValue")
	}
}
