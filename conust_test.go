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
	if string(signNegative) >= zeroOutput {
		t.Fatal("signNegative is not smaller than zeroOutput")
	}
	if string(signPositive) <= zeroOutput {
		t.Fatal("signPositive is not bigger than zeroOutput")
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

func TestDecimalSeparatorBytes(t *testing.T) {
	if positiveIntegerTerminator >= digits36[0] {
		t.Fatal("the positive decimal separator is not smaller than the digits")
	}
	if negativeIntegerTerminator <= digits36[35] {
		t.Fatal("the negative decimal separator is not smaller than the digits")
	}
}

func TestBoundaryVariables(t *testing.T) {
	if LessThanAny >= string(signNegative) {
		t.Fatal("the LessThanAny string is not smaller than the negative sign marker")
	}
	if GreaterThanAny <= string(signPositive) {
		t.Fatal("the GreaterThanAny string is not greater than the positive sign marker")
	}
}

func TestConversionI64(t *testing.T) {
	int64tests := []struct {
		name       string
		number     int64
		b10Version string
		b36Version string
	}{
		{"minimum", int64(-9223372036854775808), "4g0776627963145224191~", "4my1xazhgwxlrlr~"},
		{"bigNegativeWithTrailingZeroes", int64(-9223372030000000000), "4g077662796~", "4my1xazhk2aregv~"},
		{"divisibleBy36Negative", -1959552, "4s8040447~", "4uyt~"},
		{"mediumNegativeWitTrailingZeroes", -8000000, "4s1~", "4uv8j5r~"},
		{"minusEight", -8, "4y1~", "4yr~"},
		{"minusOne", -1, "4y8~", "4yy~"},
		{"zero", 0, "5", "5"},
		{"one", 1, "611", "611"},
		{"eight", 8, "618", "618"},
		{"mediumPositiveWithTrailingZeroes", 8000000, "678", "654rgu8"},
		{"divisibleBy36", 1959552, "671959552", "6516"},
		{"bigPositiveWithTrailingZeroes", int64(9223372030000000000), "6j922337203", "6d1y2p0ifxp8lj4"},
		{"maximum", int64(9223372036854775807), "6j9223372036854775807", "6d1y2p0ij32e8e7"},
	}

	for _, test := range int64tests {
		t.Run(test.name, func(t *testing.T) {
			enc10 := NewBase10Encoder()
			enc36 := NewBase36Encoder()
			b10 := enc10.FromInt64(test.number)
			b36 := enc36.FromInt64(test.number)

			if b10 != test.b10Version {
				t.Fatalf("B10 form for %d: '%s' instead of '%s'\n", test.number, b10, test.b10Version)
			}

			if b36 != test.b36Version {
				t.Fatalf("B36 form for %d: '%s' instead of '%s'\n", test.number, b36, test.b36Version)
			}

			var ib10 int64
			var ib36 int64

			defer func() {
				if ib10 != test.number {
					t.Fatalf("B10 form decoding error: got %d instead of %d\n", ib10, test.number)
				}

				if b36 != test.b36Version {
					t.Fatalf("B36 form decoding error: got %d instead of %d\n", ib36, test.number)
				}
			}()

			dec10 := NewBase10Decoder()
			dec36 := NewBase36Decoder()
			ib10, _ = dec10.ToInt64(b10)
			ib36, _ = dec36.ToInt64(b36)
		})
	}
}

func TestConversionI32(t *testing.T) {
	int64tests := []struct {
		name       string
		number     int32
		b10Version string
		b36Version string
	}{
		{"minimum", -2147483648, "4p7852516351~", "4t0hfz0f~"},
		{"bigNegativeWithTrailingZeroes", -2147480000, "4p785251~", "4t0hg1tr~"},
		{"divisibleBy36Negative", -1959552, "4s8040447~", "4uyt~"},
		{"mediumNegativeWitTrailingZeroes", -8000000, "4s1~", "4uv8j5r~"},
		{"minusEight", -8, "4y1~", "4yr~"},
		{"minusOne", -1, "4y8~", "4yy~"},
		{"zero", 0, "5", "5"},
		{"one", 1, "611", "611"},
		{"eight", 8, "618", "618"},
		{"mediumPositiveWithTrailingZeroes", 8000000, "678", "654rgu8"},
		{"divisibleBy36", 1959552, "671959552", "6516"},
		{"bigPositiveWithTrailingZeroes", 2147480000, "6a214748", "66zijy68"},
		{"maximum", 2147483647, "6a2147483647", "66zik0zj"},
	}

	for _, test := range int64tests {
		t.Run(test.name, func(t *testing.T) {
			enc10 := NewBase10Encoder()
			enc36 := NewBase36Encoder()
			b10 := enc10.FromInt32(test.number)
			b36 := enc36.FromInt32(test.number)

			if b10 != test.b10Version {
				t.Fatalf("B10 form for %d: '%s' instead of '%s'\n", test.number, b10, test.b10Version)
			}

			if b36 != test.b36Version {
				t.Fatalf("B36 form for %d: '%s' instead of '%s'\n", test.number, b36, test.b36Version)
			}

			var ib10 int32
			var ib36 int32

			defer func() {
				if ib10 != test.number {
					t.Fatalf("B10 form decoding error: got %d instead of %d\n", ib10, test.number)
				}

				if b36 != test.b36Version {
					t.Fatalf("B36 form decoding error: got %d instead of %d\n", ib36, test.number)
				}
			}()

			dec10 := NewBase10Decoder()
			dec36 := NewBase36Decoder()
			ib10, _ = dec10.ToInt32(b10)
			ib36, _ = dec36.ToInt32(b36)
		})
	}
}

func TestAnalysis(t *testing.T) {
	analysisTests := []struct {
		name     string
		input    string
		ok       bool
		empty    bool
		zero     bool
		positive bool
		intNZ    string
		intLen   int
		fracLZ   int
		fracNZ   string
	}{
		{name: "bad input", input: "asdf", ok: false},
		{name: "epmty", input: "", ok: true, empty: true},
		{name: "minus only", input: "-", ok: false},
		{name: "plus only", input: "-", ok: false},
		{name: "minus zero", input: "-0", ok: true, empty: false, zero: true},
		{name: "plus zero", input: "+0", ok: true, empty: false, zero: true},
		{name: "zero", input: "0", ok: true, empty: false, zero: true},
		{name: "more zeroes", input: "000", ok: true, empty: false, zero: true},
		{name: "frac zeroes", input: "000.000", ok: true, empty: false, zero: true},
		{name: "one", input: "1", ok: true, empty: false, zero: false, positive: true, intNZ: "1", intLen: 1},
		{name: "nine", input: "9", ok: true, empty: false, zero: false, positive: true, intNZ: "9", intLen: 1},
		{name: "single int and frac 1", input: "1.9", ok: true, empty: false, zero: false, positive: true, intNZ: "1", intLen: 1, fracLZ: 0, fracNZ: "9"},
		{name: "single int and frac 2", input: "9.1", ok: true, empty: false, zero: false, positive: true, intNZ: "9", intLen: 1, fracLZ: 0, fracNZ: "1"},
		{name: "leading zeroes int", input: "00100", ok: true, empty: false, zero: false, positive: true, intNZ: "1", intLen: 3, fracLZ: 0, fracNZ: ""},
		{name: "all the things", input: "02900.00410", ok: true, empty: false, zero: false, positive: true, intNZ: "29", intLen: 4, fracLZ: 2, fracNZ: "41"},
		{name: "neagtive things", input: "-02900.00410", ok: true, empty: false, zero: false, positive: false, intNZ: "29", intLen: 4, fracLZ: 2, fracNZ: "41"},
	}

	enc := (NewBase10Encoder()).(*Base10Encoder)
	for _, i := range analysisTests {
		t.Run(i.name, func(t *testing.T) {
			enc.Encode(i.input)
			// fmt.Println("Encoded", i.input, "to", out)

			if i.ok != enc.ok {
				t.Fatalf("OK expected: %v, got %v\n", i.ok, enc.ok)
			}
			if !i.ok {
				return
			}

			if i.empty != enc.empty {
				t.Fatalf("isEmpty expected: %v, got %v\n", i.empty, enc.empty)
			}
			if i.empty {
				return
			}

			if i.zero != enc.zero {
				t.Fatalf("isZero expected: %v, got %v\n", i.zero, enc.zero)
			}
			if i.zero {
				return
			}

			if i.positive != enc.positive {
				t.Fatalf("isZero expected: %v, got %v\n", i.positive, enc.positive)
			}

			tempIntNZ := i.input[enc.intNonZeroFrom:enc.intNonZeroTo]
			if tempIntNZ != i.intNZ {
				t.Fatalf("int NZ expected: %v, got %v (in %q [%d : %d]\n", i.intNZ, tempIntNZ, i.input, enc.intNonZeroFrom, enc.intNonZeroTo)
			}

			tempIntLen := enc.intEnd - enc.intNonZeroFrom + 1
			if i.intLen != tempIntLen {
				t.Fatalf("int length expected: %d, got %d\n", i.intLen, tempIntLen)
			}

			fracNZ := i.input[enc.fracNonZeroFrom:enc.fracNonZeroTo]
			if i.fracNZ != fracNZ {
				t.Fatalf("empty frac NZ expected: %s, got: %s (in %s [%d : %d])", i.fracNZ, fracNZ, i.input, enc.fracNonZeroFrom, enc.fracNonZeroTo)
			}

			if i.fracLZ != enc.fracLeadingZeroCount {
				t.Fatalf("frac leadig zero count expected: %d, got%d\n", i.fracLZ, enc.fracLeadingZeroCount)
			}
		})
	}
}
