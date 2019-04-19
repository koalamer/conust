package conust

import "testing"

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
		{name: "epmty", input: "", ok: true, empty: true},

		{name: "minus only", input: "-", ok: false},
		{name: "plus only", input: "-", ok: false},
		{name: "smaller non digit", input: "!", ok: false},
		{name: "inbetween non digit", input: "=", ok: false},
		{name: "greater non digit", input: "}", ok: false},

		{name: "all numeric digits", input: "1234567890", ok: true, empty: false, zero: false, positive: true, intNZ: "123456789", intLen: 10},
		{name: "all alpha digits", input: "abcdefghijklmnopqrstuvwxyz", ok: true, empty: false, zero: false, positive: true, intNZ: "abcdefghijklmnopqrstuvwxyz", intLen: 26},

		{name: "minus zero", input: "-0", ok: true, empty: false, zero: true},
		{name: "plus zero", input: "+0", ok: true, empty: false, zero: true},
		{name: "zero", input: "0", ok: true, empty: false, zero: true},
		{name: "more zeroes", input: "000", ok: true, empty: false, zero: true},
		{name: "frac zeroes", input: "000.000", ok: true, empty: false, zero: true},

		{name: "one", input: "1", ok: true, empty: false, zero: false, positive: true, intNZ: "1", intLen: 1},
		{name: "nine", input: "9", ok: true, empty: false, zero: false, positive: true, intNZ: "9", intLen: 1},
		{name: "nine", input: "a", ok: true, empty: false, zero: false, positive: true, intNZ: "a", intLen: 1},
		{name: "nine", input: "z", ok: true, empty: false, zero: false, positive: true, intNZ: "z", intLen: 1},

		{name: "single int and frac 1", input: "1.9", ok: true, empty: false, zero: false, positive: true, intNZ: "1", intLen: 1, fracLZ: 0, fracNZ: "9"},
		{name: "single int and frac 2", input: "9.1", ok: true, empty: false, zero: false, positive: true, intNZ: "9", intLen: 1, fracLZ: 0, fracNZ: "1"},
		{name: "leading zeroes int", input: "00100", ok: true, empty: false, zero: false, positive: true, intNZ: "1", intLen: 3, fracLZ: 0, fracNZ: ""},
		{name: "all the things", input: "02900.00410", ok: true, empty: false, zero: false, positive: true, intNZ: "29", intLen: 4, fracLZ: 2, fracNZ: "41"},
		{name: "neagtive all the things", input: "-02900.00410", ok: true, empty: false, zero: false, positive: false, intNZ: "29", intLen: 4, fracLZ: 2, fracNZ: "41"},
		{name: "base 36 general", input: "-ace00.00decade", ok: true, empty: false, zero: false, positive: false, intNZ: "ace", intLen: 5, fracLZ: 2, fracNZ: "decade"},
		{name: "negative base 36 general", input: "-0can00.00000do0", ok: true, empty: false, zero: false, positive: false, intNZ: "can", intLen: 5, fracLZ: 5, fracNZ: "do"},
		{name: "maximum int length", input: "-0huge0000000000000000000000000000000.00410", ok: true, empty: false, zero: false, positive: false, intNZ: "huge", intLen: 35, fracLZ: 2, fracNZ: "41"},
		{name: "maximum frac leading zero count", input: "0.00000000000000000000000000000000000410", ok: true, empty: false, zero: false, positive: true, intNZ: "", intLen: 0, fracLZ: 35, fracNZ: "41"},
		{name: "too long integer part", input: "-0overkill0000000000000000000000000000000.00410", ok: false},
		{name: "too many frac leading zeros", input: "0.000000000000000000000000000000000000410", ok: false},
	}

	enc := NewCodec()
	for _, i := range analysisTests {
		t.Run(i.name, func(t *testing.T) {
			_, ok := enc.Encode(i.input)
			encoderState := enc.(*codec)

			if i.ok != ok {
				t.Fatalf("OK output expected: %v, got %v\n", i.ok, ok)
			}
			if i.ok != encoderState.ok {
				t.Fatalf("OK inner state expected: %v, got %v\n", i.ok, encoderState.ok)
			}
			if !i.ok {
				return
			}

			if i.empty != encoderState.empty {
				t.Fatalf("Empty expected: %v, got %v\n", i.empty, encoderState.empty)
			}
			if i.empty {
				return
			}

			if i.zero != encoderState.zero {
				t.Fatalf("Zero expected: %v, got %v\n", i.zero, encoderState.zero)
			}
			if i.zero {
				return
			}

			if i.positive != encoderState.positive {
				t.Fatalf("Positive expected: %v, got %v\n", i.positive, encoderState.positive)
			}

			tempIntNZ := i.input[encoderState.intNonZeroFrom:encoderState.intNonZeroTo]
			if tempIntNZ != i.intNZ {
				t.Fatalf("int non zero part expected: %v, got %v (in %q [%d : %d]\n", i.intNZ, tempIntNZ, i.input, encoderState.intNonZeroFrom, encoderState.intNonZeroTo)
			}

			tempIntLen := encoderState.intTo - encoderState.intNonZeroFrom
			if i.intLen != tempIntLen {
				t.Fatalf("int length expected: %d, got %d\n", i.intLen, tempIntLen)
			}

			fracNZ := i.input[encoderState.fracNonZeroFrom:encoderState.fracNonZeroTo]
			if i.fracNZ != fracNZ {
				t.Fatalf("empty frac non zero part expected: %s, got: %s (in %s [%d : %d])", i.fracNZ, fracNZ, i.input, encoderState.fracNonZeroFrom, encoderState.fracNonZeroTo)
			}

			if i.fracLZ != encoderState.fracLeadingZeroCount {
				t.Fatalf("frac leadig zero count expected: %d, got%d\n", i.fracLZ, encoderState.fracLeadingZeroCount)
			}
		})
	}
}
