package conust

import (
	"testing"
)

func TestEncodingAnalysis(t *testing.T) {
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
		output   string
	}{
		{name: "epmty", input: "", ok: true, empty: true},

		{name: "minus only", input: "-", ok: false},
		{name: "plus only", input: "-", ok: false},
		{name: "smaller non digit", input: "!", ok: false},
		{name: "inbetween non digit", input: "=", ok: false},
		{name: "greater non digit", input: "}", ok: false},

		{name: "all numeric digits", input: "1234567890", ok: true, empty: false, zero: false, positive: true, intNZ: "123456789", intLen: 10, output: "6a123456789"},
		{name: "all alpha digits", input: "abcdefghijklmnopqrstuvwxyz", ok: true, empty: false, zero: false, positive: true, intNZ: "abcdefghijklmnopqrstuvwxyz", intLen: 26, output: "6qabcdefghijklmnopqrstuvwxyz"},

		{name: "minus zero", input: "-0", ok: true, empty: false, zero: true, output: "0"},
		{name: "plus zero", input: "+0", ok: true, empty: false, zero: true, output: "0"},
		{name: "zero", input: "0", ok: true, empty: false, zero: true, output: "0"},
		{name: "more zeroes", input: "000", ok: true, empty: false, zero: true, output: "0"},
		{name: "frac zeroes", input: "000.000", ok: true, empty: false, zero: true, output: "0"},

		{name: "one", input: "1", ok: true, empty: false, zero: false, positive: true, intNZ: "1", intLen: 1, output: "611"},
		{name: "nine", input: "9", ok: true, empty: false, zero: false, positive: true, intNZ: "9", intLen: 1, output: "619"},
		{name: "nine", input: "a", ok: true, empty: false, zero: false, positive: true, intNZ: "a", intLen: 1, output: "61a"},
		{name: "nine", input: "z", ok: true, empty: false, zero: false, positive: true, intNZ: "z", intLen: 1, output: "61z"},

		{name: "single int and frac 1", input: "1.9", ok: true, empty: false, zero: false, positive: true, intNZ: "1", intLen: 1, fracLZ: 0, fracNZ: "9", output: "611.z9"},
		{name: "single int and frac 2", input: "9.1", ok: true, empty: false, zero: false, positive: true, intNZ: "9", intLen: 1, fracLZ: 0, fracNZ: "1", output: "619.z1"},
		{name: "leading zeroes int", input: "00100", ok: true, empty: false, zero: false, positive: true, intNZ: "1", intLen: 3, fracLZ: 0, fracNZ: "", output: "631"},
		{name: "all the things", input: "02900.00410", ok: true, empty: false, zero: false, positive: true, intNZ: "29", intLen: 4, fracLZ: 2, fracNZ: "41", output: "6429.x41"},
		{name: "neagtive all the things", input: "-02900.00410", ok: true, empty: false, zero: false, positive: false, intNZ: "29", intLen: 4, fracLZ: 2, fracNZ: "41", output: "4vxq~2vy~"},
		{name: "base 36 general", input: "-ace00.00decade", ok: true, empty: false, zero: false, positive: false, intNZ: "ace", intLen: 5, fracLZ: 2, fracNZ: "decade", output: "4upnl~2mlnpml~"},
		{name: "negative base 36 general", input: "-0can00.00000do0", ok: true, empty: false, zero: false, positive: false, intNZ: "can", intLen: 5, fracLZ: 5, fracNZ: "do", output: "4unpc~5mb~"},
		{name: "maximum int length", input: "-0huge0000000000000000000000000000000.00410", ok: true, empty: false, zero: false, positive: false, intNZ: "huge", intLen: 35, fracLZ: 2, fracNZ: "41", output: "40i5jl~2vy~"},
		{name: "maximum frac leading zero count", input: "0.00000000000000000000000000000000000410", ok: true, empty: false, zero: false, positive: true, intNZ: "", intLen: 0, fracLZ: 35, fracNZ: "41", output: "610.041"},
		{name: "too long integer part", input: "-0overkill0000000000000000000000000000.00410", ok: false, output: ""},
		{name: "too many frac leading zeros", input: "0.000000000000000000000000000000000000410", ok: false, output: ""},
	}

	enc := NewCodec()
	for _, i := range analysisTests {
		t.Run(i.name, func(t *testing.T) {
			output, ok := enc.Encode(i.input)
			codecState := enc.(*codec)

			if i.ok != ok {
				t.Fatalf("OK output expected: %v, got %v\n", i.ok, ok)
			}
			if i.ok != codecState.ok {
				t.Fatalf("OK inner state expected: %v, got %v\n", i.ok, codecState.ok)
			}
			if !i.ok {
				return
			}

			if i.empty != codecState.empty {
				t.Fatalf("Empty expected: %v, got %v\n", i.empty, codecState.empty)
			}
			if i.empty {
				return
			}

			if i.zero != codecState.zero {
				t.Fatalf("Zero expected: %v, got %v\n", i.zero, codecState.zero)
			}
			if i.zero {
				return
			}

			if i.positive != codecState.positive {
				t.Fatalf("Positive expected: %v, got %v\n", i.positive, codecState.positive)
			}

			tempIntNZ := i.input[codecState.intNonZeroFrom:codecState.intNonZeroTo]
			if tempIntNZ != i.intNZ {
				t.Fatalf("int non zero part expected: %v, got %v (in %q [%d : %d]\n", i.intNZ, tempIntNZ, i.input, codecState.intNonZeroFrom, codecState.intNonZeroTo)
			}

			tempIntLen := codecState.intTo - codecState.intNonZeroFrom
			if i.intLen != tempIntLen {
				t.Fatalf("int length expected: %d, got %d\n", i.intLen, tempIntLen)
			}

			fracNZ := i.input[codecState.fracNonZeroFrom:codecState.fracNonZeroTo]
			if i.fracNZ != fracNZ {
				t.Fatalf("empty frac non zero part expected: %s, got: %s (in %s [%d : %d])", i.fracNZ, fracNZ, i.input, codecState.fracNonZeroFrom, codecState.fracNonZeroTo)
			}

			if i.fracLZ != codecState.fracLeadingZeroCount {
				t.Fatalf("frac leadig zero count expected: %d, got %d\n", i.fracLZ, codecState.fracLeadingZeroCount)
			}

			if i.output != output {
				t.Fatalf("output expected: %s, got %s\n", i.output, output)
			}
		})
	}
}

func TestDecodingAnalysis(t *testing.T) {
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
		output   string
	}{
		{name: "epmty", input: "", ok: true, empty: true, output: ""},

		{name: "zerro", input: "5", ok: true, zero: true, empty: false, output: "0"},

		{name: "one", input: "611", ok: true, empty: false, zero: false, positive: true, intLen: 1, intNZ: "1", output: "1"},
		{name: "zf", input: "62zf", ok: true, empty: false, zero: false, positive: true, intLen: 2, intNZ: "zf", output: "zf"},
		{name: "trimmed", input: "6666", ok: true, empty: false, zero: false, positive: true, intLen: 6, intNZ: "66", output: "660000"},
		{name: "negative one", input: "4yy~", ok: true, empty: false, zero: false, positive: false, intLen: 1, intNZ: "y", output: "-1"},
		{name: "negative trimmed", input: "4ttt~", ok: true, empty: false, zero: false, positive: false, intLen: 6, intNZ: "tt", output: "-660000"},

		{name: "fractional", input: "611.z9", ok: true, empty: false, zero: false, positive: true, intLen: 1, intNZ: "1", fracLZ: 0, fracNZ: "9", output: "1.9"},
		{name: "fractional 2", input: "6429.x41", ok: true, empty: false, zero: false, positive: true, intLen: 4, intNZ: "29", fracLZ: 2, fracNZ: "41", output: "2900.0041"},
		{name: "fractional negative", input: "4vxq~2vy~", ok: true, empty: false, zero: false, positive: false, intLen: 4, intNZ: "xq", fracLZ: 2, fracNZ: "vy", output: "-2900.0041"},
		{name: "fractional negative 2", input: "4upnl~2mlnpml~", ok: true, empty: false, zero: false, positive: false, intLen: 5, intNZ: "pnl", fracLZ: 2, fracNZ: "mlnpml", output: "-ace00.00decade"},

		{name: "maximum int length", input: "40i5jl~2vy~", ok: true, empty: false, zero: false, positive: false, intLen: 35, intNZ: "i5jl", fracLZ: 2, fracNZ: "vy", output: "-huge0000000000000000000000000000000.0041"},
		{name: "maximum frac leading zeros", input: "610.041", ok: true, empty: false, zero: false, positive: true, intLen: 1, intNZ: "0", fracLZ: 35, fracNZ: "41", output: "0.0000000000000000000000000000000000041"},
	}
	dec := NewCodec()
	for _, i := range analysisTests {
		t.Run(i.name, func(t *testing.T) {
			output, ok := dec.Decode(i.input)
			codecState := dec.(*codec)

			if i.ok != ok {
				t.Fatalf("OK output expected: %v, got %v\n", i.ok, ok)
			}
			if i.ok != codecState.ok {
				t.Fatalf("OK inner state expected: %v, got %v\n", i.ok, codecState.ok)
			}
			if !i.ok {
				return
			}

			if i.empty != codecState.empty {
				t.Fatalf("Empty expected: %v, got %v\n", i.empty, codecState.empty)
			}
			if i.empty {
				return
			}

			if i.zero != codecState.zero {
				t.Fatalf("Zero expected: %v, got %v\n", i.zero, codecState.zero)
			}
			if i.zero {
				return
			}

			if i.positive != codecState.positive {
				t.Fatalf("Positive expected: %v, got %v\n", i.positive, codecState.positive)
			}

			tempIntNZ := i.input[codecState.intNonZeroFrom:codecState.intNonZeroTo]
			if tempIntNZ != i.intNZ {
				t.Fatalf("int non zero part expected: %v, got %v (in %q [%d : %d]\n", i.intNZ, tempIntNZ, i.input, codecState.intNonZeroFrom, codecState.intNonZeroTo)
			}

			tempIntLen := codecState.intTo - codecState.intNonZeroFrom
			if i.intLen != tempIntLen {
				t.Fatalf("int length expected: %d, got %d\n", i.intLen, tempIntLen)
			}

			fracNZ := i.input[codecState.fracNonZeroFrom:codecState.fracNonZeroTo]
			if i.fracNZ != fracNZ {
				t.Fatalf("empty frac non zero part expected: %s, got: %s (in %s [%d : %d])", i.fracNZ, fracNZ, i.input, codecState.fracNonZeroFrom, codecState.fracNonZeroTo)
			}

			if i.fracLZ != codecState.fracLeadingZeroCount {
				t.Fatalf("frac leadig zero count expected: %d, got %d\n", i.fracLZ, codecState.fracLeadingZeroCount)
			}

			if i.output != output {
				t.Fatalf("output expected: %s, got %s\n", i.output, output)
			}
		})
	}
}
func TestCodec(t *testing.T) {
	codecTests := []struct {
		name    string
		input   string
		encoded string
		decoded string
	}{
		{name: "empty", input: "", encoded: "", decoded: ""},

		{name: "zero 1", input: "0", encoded: "5", decoded: "0"},
		{name: "zero 2", input: "+000", encoded: "5", decoded: "0"},
		{name: "zero 3", input: "-000", encoded: "5", decoded: "0"},
		{name: "zero 4", input: "000.0000", encoded: "5", decoded: "0"},

		{name: "all digits", input: "1234567890abcdefghij.klmnopqrstuvwxyz", encoded: "6k1234567890abcdefghij.zklmnopqrstuvwxyz", decoded: "1234567890abcdefghij.klmnopqrstuvwxyz"},
		{name: "negative all digits", input: "-1234567890abcdefghij.klmnopqrstuvwxyz", encoded: "4fyxwvutsrqzponmlkjihg~0fedcba9876543210~", decoded: "-1234567890abcdefghij.klmnopqrstuvwxyz"},

		{name: "one", input: "1", encoded: "611", decoded: "1"},
		{name: "ugly one", input: "+00001", encoded: "611", decoded: "1"},
		{name: "negative one", input: "-1", encoded: "4yy~", decoded: "-1"},
		{name: "ugly negative one", input: "-000001", encoded: "4yy~", decoded: "-1"},
		{name: "ugly positive int", input: "+00000123000", encoded: "66123", decoded: "123000"},
		{name: "ugly negative int", input: "-00000123000", encoded: "4tyxw~", decoded: "-123000"},
		{name: "fractional", input: "54321.12345", encoded: "6554321.z12345", decoded: "54321.12345"},
		{name: "negative fractional", input: "-54321.12345", encoded: "4uuvwxy~0yxwvu~", decoded: "-54321.12345"},
		{name: "ugly fractional", input: "+00054321000.00012345000", encoded: "6854321.w12345", decoded: "54321000.00012345"},
		{name: "ugly negative fractional", input: "-00054321000.00012345000", encoded: "4ruvwxy~3yxwvu~", decoded: "-54321000.00012345"},
		{name: "cowboy hat", input: "cowboy.hat", encoded: "66cowboy.zhat", decoded: "cowboy.hat"},
		{name: "negative cowboy hat", input: "-cowboy.hat", encoded: "4tnb3ob1~0ip6~", decoded: "-cowboy.hat"},
		{name: "maximum int length", input: "12345678901234567890123456789012345.1", encoded: "6z12345678901234567890123456789012345.z1", decoded: "12345678901234567890123456789012345.1"},
		{name: "maximum negative int length", input: "-12345678901234567890123456789012345.1", encoded: "40yxwvutsrqzyxwvutsrqzyxwvutsrqzyxwvu~0y~", decoded: "-12345678901234567890123456789012345.1"},
		{name: "maximum fracleading zero count", input: "3.000000000000000000000000000000000004325430", encoded: "613.0432543", decoded: "3.00000000000000000000000000000000000432543"},
	}
	codec := NewCodec()
	for _, i := range codecTests {
		t.Run(i.name, func(t *testing.T) {
			encoded, _ := codec.Encode(i.input)

			if i.encoded != encoded {
				t.Fatalf("Encoding expected: %v, got %v\n", i.encoded, encoded)
			}

			decoded, _ := codec.Decode(encoded)
			if i.decoded != decoded {
				t.Fatalf("Decoding expected: %v, got %v\n", i.decoded, decoded)
			}
		})
	}
}
