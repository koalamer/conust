package conust

import (
	"testing"
)

func TestArrayReversion(t *testing.T) {
	tests := []struct {
		name     string
		forward  []byte
		backward []byte
	}{
		{"base10", digits36[0:10], digits10Reversed[:]},
		{"base36", digits36[:], digits36Reversed[:]},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			offset := len(testCase.forward) - 1
			if len(testCase.forward) != len(testCase.backward) {
				t.Fatal("Forward and backward digits arrays are of different length")
			}

			for i, c := range testCase.backward {
				if c != testCase.forward[offset-i] {
					t.Fatalf("digits %d backward[%d] = %s but forward[%d] = %s",
						len(testCase.forward),
						i,
						string(testCase.backward[i]),
						offset-i,
						string(testCase.forward[offset-i]),
					)
				}
			}
		})
	}
}

func TestArraySortedness(t *testing.T) {
	tests := []struct {
		name      string
		arr       []byte
		ascending bool
	}{
		{"digits10Reversed", digits10Reversed[:], false},
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
	if string(signNegative10) >= zeroOutput {
		t.Fatal("signNegative10 is not smaller than zeroOutput")
	}
	if signNegative10%2 != 0 {
		t.Fatal("signNegative10 is not even")
	}

	if string(signPositive10) <= zeroOutput {
		t.Fatal("signPositive10 is not bigger than zeroOutput")
	}
	if signPositive10%2 != 0 {
		t.Fatal("signPositive10 is not even")
	}

	if string(signNegative36) >= zeroOutput {
		t.Fatal("signNegative36 is not smaller than zeroOutput")
	}
	if signNegative36%2 != 1 {
		t.Fatal("signNegative36 is not odd")
	}

	if string(signPositive36) <= zeroOutput {
		t.Fatal("signPositive36 is not bigger than zeroOutput")
	}
	if signPositive36%2 != 1 {
		t.Fatal("signPositive36 is not odd")
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

	if len(trailing0) != 1 {
		t.Fatal("trailing0 is not one character long")
	}
	if trailing0[0] != digits36[0] {
		t.Fatal("trailing0 is not the [0] digit")
	}
}
