package conust

import (
	"fmt"
	"testing"
)

func TestEncodingAnalysis(t *testing.T) {
	analysisTests := []struct {
		name              string
		input             string
		ok                bool
		empty             bool
		zero              bool
		positive          bool
		intNZ             string
		intLen            int
		fracLZ            int
		fracNZ            string
		output            string
		magnitude         int
		magnitudePositive bool
	}{
		{name: "epmty", input: "", ok: true, empty: true},

		{name: "minus only", input: "-", ok: false},
		{name: "plus only", input: "-", ok: false},
		{name: "smaller non digit", input: "!", ok: false},
		{name: "inbetween non digit", input: "=", ok: false},
		{name: "greater non digit", input: "}", ok: false},

		{name: "all numeric digits", input: "1234567890", ok: true, empty: false, zero: false, positive: true, intNZ: "123456789", intLen: 10, output: "7a123456789", magnitude: 10, magnitudePositive: true},
		{name: "all alpha digits", input: "abcdefghijklmnopqrstuvwxyz", ok: true, empty: false, zero: false, positive: true, intNZ: "abcdefghijklmnopqrstuvwxyz", intLen: 26, output: "7qabcdefghijklmnopqrstuvwxyz", magnitude: 26, magnitudePositive: true},

		{name: "minus zero", input: "-0", ok: true, empty: false, zero: true, output: "5"},
		{name: "plus zero", input: "+0", ok: true, empty: false, zero: true, output: "5"},
		{name: "zero", input: "0", ok: true, empty: false, zero: true, output: "5"},
		{name: "more zeroes", input: "000", ok: true, empty: false, zero: true, output: "5"},
		{name: "frac zeroes", input: "000.000", ok: true, empty: false, zero: true, output: "5"},

		{name: "one", input: "1", ok: true, empty: false, zero: false, positive: true, intNZ: "1", intLen: 1, output: "711", magnitude: 1, magnitudePositive: true},
		{name: "nine", input: "9", ok: true, empty: false, zero: false, positive: true, intNZ: "9", intLen: 1, output: "719", magnitude: 1, magnitudePositive: true},
		{name: "nine", input: "a", ok: true, empty: false, zero: false, positive: true, intNZ: "a", intLen: 1, output: "71a", magnitude: 1, magnitudePositive: true},
		{name: "nine", input: "z", ok: true, empty: false, zero: false, positive: true, intNZ: "z", intLen: 1, output: "71z", magnitude: 1, magnitudePositive: true},

		{name: "frac 1", input: "0.9", ok: true, empty: false, zero: false, positive: true, intNZ: "", intLen: 0, fracLZ: 0, fracNZ: "9", output: "6z9", magnitude: 0, magnitudePositive: false},
		{name: "frac 2", input: "0.000123", ok: true, empty: false, zero: false, positive: true, intNZ: "", intLen: 0, fracLZ: 3, fracNZ: "123", output: "6w123", magnitude: 3, magnitudePositive: false},
		{name: "frac 3", input: "0.000000000000000000000000000000000000000000000000000000000000000000000000end",
			ok: true, empty: false, zero: false, positive: true, intNZ: "", intLen: 0, fracLZ: 72, fracNZ: "end", output: "600vend", magnitude: 72, magnitudePositive: false},

		{name: "negative frac 1", input: "-0.9", ok: true, empty: false, zero: false, positive: false, intNZ: "", intLen: 0, fracLZ: 0, fracNZ: "9", output: "40q~", magnitude: 0, magnitudePositive: false},
		{name: "negative frac 2", input: "-0.000123", ok: true, empty: false, zero: false, positive: false, intNZ: "", intLen: 0, fracLZ: 3, fracNZ: "123", output: "43yxw~", magnitude: 3, magnitudePositive: false},
		{name: "negative frac 3", input: "-0.000000000000000000000000000000000000000000000000000000000000000000000000end",
			ok: true, empty: false, zero: false, positive: false, intNZ: "", intLen: 0, fracLZ: 72, fracNZ: "end", output: "4zz4lcm~", magnitude: 72, magnitudePositive: false},

		{name: "single int and frac 1", input: "1.9", ok: true, empty: false, zero: false, positive: true, intNZ: "1", intLen: 1, fracLZ: 0, fracNZ: "9", output: "7119", magnitude: 1, magnitudePositive: true},
		{name: "single int and frac 2", input: "9.1", ok: true, empty: false, zero: false, positive: true, intNZ: "9", intLen: 1, fracLZ: 0, fracNZ: "1", output: "7191", magnitude: 1, magnitudePositive: true},
		{name: "leading zeroes int", input: "00100", ok: true, empty: false, zero: false, positive: true, intNZ: "1", intLen: 3, fracLZ: 0, fracNZ: "", output: "731", magnitude: 3, magnitudePositive: true},
		{name: "all the things", input: "02900.00410", ok: true, empty: false, zero: false, positive: true, intNZ: "29", intLen: 4, fracLZ: 2, fracNZ: "41", output: "7429000041", magnitude: 4, magnitudePositive: true},
		{name: "neagtive all the things", input: "-02900.00410", ok: true, empty: false, zero: false, positive: false, intNZ: "29", intLen: 4, fracLZ: 2, fracNZ: "41", output: "3vxqzzzzvy~", magnitude: 4, magnitudePositive: true},
		{name: "base 36 general", input: "-ace00.00decade", ok: true, empty: false, zero: false, positive: false, intNZ: "ace", intLen: 5, fracLZ: 2, fracNZ: "decade", output: "3upnlzzzzmlnpml~", magnitude: 5, magnitudePositive: true},
		{name: "negative base 36 general", input: "-0can00.00000do0", ok: true, empty: false, zero: false, positive: false, intNZ: "can", intLen: 5, fracLZ: 5, fracNZ: "do", output: "3unpczzzzzzzmb~", magnitude: 5, magnitudePositive: true},
		{name: "big negative number", input: "-huge000000000000000000000000000000000000000000000.00", ok: true, empty: false, zero: false, positive: false, intNZ: "huge", intLen: 49, fracLZ: 0, fracNZ: "", output: "30ki5jl~", magnitude: 49, magnitudePositive: true},
		{name: "big number", input: "huge000000000000000000000000000000000000000000000.00", ok: true, empty: false, zero: false, positive: true, intNZ: "huge", intLen: 49, fracLZ: 0, fracNZ: "", output: "7zfhuge", magnitude: 49, magnitudePositive: true},
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

			tempIntNZ := i.input[codecState.intSignificantFrom:codecState.intSignificantTo]
			if tempIntNZ != i.intNZ {
				t.Fatalf("int non zero part expected: %v, got %v (in %q [%d : %d]\n", i.intNZ, tempIntNZ, i.input, codecState.intSignificantFrom, codecState.intSignificantTo)
			}

			tempIntLen := codecState.intTo - codecState.intSignificantFrom
			if i.intLen != tempIntLen {
				t.Fatalf("int length expected: %d, got %d\n", i.intLen, tempIntLen)
			}

			fracNZ := i.input[codecState.fracSignificantFrom:codecState.fracSignificantTo]
			if i.fracNZ != fracNZ {
				t.Fatalf("empty frac non zero part expected: %s, got: %s (in %s [%d : %d])", i.fracNZ, fracNZ, i.input, codecState.fracSignificantFrom, codecState.fracSignificantTo)
			}

			if i.fracLZ != codecState.fracLeadingZeroCount {
				t.Fatalf("frac leadig zero count expected: %d, got %d\n", i.fracLZ, codecState.fracLeadingZeroCount)
			}

			if i.magnitude != codecState.magnitude {
				t.Fatalf("magnitude expected: %d, got %d\n", i.magnitude, codecState.magnitude)
			}

			if i.magnitudePositive != codecState.magnitudePositive {
				t.Fatalf("magnitudePositive expected: %v, got %v\n", i.magnitudePositive, codecState.magnitudePositive)
			}

			if i.output != output {
				t.Fatalf("output expected: %s, got %s\n", i.output, output)
			}
		})
	}
}

func TestDecodingAnalysis(t *testing.T) {
	analysisTests := []struct {
		name              string
		input             string
		ok                bool
		empty             bool
		zero              bool
		positive          bool
		intNZ             string
		intLen            int
		fracLZ            int
		fracNZ            string
		output            string
		magnitude         int
		magnitudePositive bool
	}{
		{name: "epmty", input: "", ok: true, empty: true, output: ""},

		{name: "zerro", input: "5", ok: true, zero: true, empty: false, output: "0"},

		{name: "one", input: "711", ok: true, empty: false, zero: false, positive: true, intLen: 1, intNZ: "1", output: "1", magnitude: 1, magnitudePositive: true},
		{name: "zf", input: "72zf", ok: true, empty: false, zero: false, positive: true, intLen: 2, intNZ: "zf", output: "zf", magnitude: 2, magnitudePositive: true},
		{name: "trimmed", input: "7666", ok: true, empty: false, zero: false, positive: true, intLen: 6, intNZ: "66", output: "660000", magnitude: 6, magnitudePositive: true},
		{name: "negative one", input: "3yy~", ok: true, empty: false, zero: false, positive: false, intLen: 1, intNZ: "y", output: "-1", magnitude: 1, magnitudePositive: true},
		{name: "negative trimmed", input: "3ttt~", ok: true, empty: false, zero: false, positive: false, intLen: 6, intNZ: "tt", output: "-660000", magnitude: 6, magnitudePositive: true},

		{name: "fractional", input: "7119", ok: true, empty: false, zero: false, positive: true, intLen: 1, intNZ: "1", fracLZ: 0, fracNZ: "9", output: "1.9", magnitude: 1, magnitudePositive: true},
		{name: "fractional 2", input: "7429000041", ok: true, empty: false, zero: false, positive: true, intLen: 4, intNZ: "2900", fracLZ: 0, fracNZ: "0041", output: "2900.0041", magnitude: 4, magnitudePositive: true},
		{name: "fractional negative", input: "3vxqzzzzvy~", ok: true, empty: false, zero: false, positive: false, intLen: 4, intNZ: "xqzz", fracLZ: 0, fracNZ: "zzvy", output: "-2900.0041", magnitude: 4, magnitudePositive: true},
		{name: "fractional negative 2", input: "3upnlzzzzmlnpml~", ok: true, empty: false, zero: false, positive: false, intLen: 5, intNZ: "pnlzz", fracLZ: 0, fracNZ: "zzmlnpml", output: "-ace00.00decade", magnitude: 5, magnitudePositive: true},

		{name: "maximum int length", input: "30yi5jlzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzvy~", ok: true, empty: false, zero: false, positive: false, intLen: 35, intNZ: "i5jlzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", fracLZ: 0, fracNZ: "zzvy", output: "-huge0000000000000000000000000000000.0041", magnitude: 35, magnitudePositive: true},
		{name: "maximum frac leading zeros", input: "60y41", ok: true, empty: false, zero: false, positive: true, intLen: 0, intNZ: "", fracLZ: 35, fracNZ: "41", output: "0.0000000000000000000000000000000000041", magnitude: 35, magnitudePositive: false},
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

			tempIntNZ := i.input[codecState.intSignificantFrom:codecState.intSignificantTo]
			if tempIntNZ != i.intNZ {
				t.Fatalf("int non zero part expected: %v, got %v (in %q [%d : %d]\n", i.intNZ, tempIntNZ, i.input, codecState.intSignificantFrom, codecState.intSignificantTo)
			}

			tempIntLen := codecState.intTo - codecState.intSignificantFrom
			if i.intLen != tempIntLen {
				t.Fatalf("int length expected: %d, got %d\n", i.intLen, tempIntLen)
			}

			fracNZ := i.input[codecState.fracSignificantFrom:codecState.fracSignificantTo]
			if i.fracNZ != fracNZ {
				t.Fatalf("empty frac non zero part expected: %s, got: %s (in %s [%d : %d])", i.fracNZ, fracNZ, i.input, codecState.fracSignificantFrom, codecState.fracSignificantTo)
			}

			if i.fracLZ != codecState.fracLeadingZeroCount {
				t.Fatalf("frac leadig zero count expected: %d, got %d\n", i.fracLZ, codecState.fracLeadingZeroCount)
			}

			if i.magnitude != codecState.magnitude {
				t.Fatalf("magnitude expected: %d, got %d\n", i.magnitude, codecState.magnitude)
			}

			if i.magnitudePositive != codecState.magnitudePositive {
				t.Fatalf("magnitudePositive expected: %v, got %v\n", i.magnitudePositive, codecState.magnitudePositive)
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

		{name: "all digits", input: "1234567890abcdefghij.klmnopqrstuvwxyz", encoded: "7k1234567890abcdefghijklmnopqrstuvwxyz", decoded: "1234567890abcdefghij.klmnopqrstuvwxyz"},
		{name: "negative all digits", input: "-1234567890abcdefghij.klmnopqrstuvwxyz", encoded: "3fyxwvutsrqzponmlkjihgfedcba9876543210~", decoded: "-1234567890abcdefghij.klmnopqrstuvwxyz"},

		{name: "holes in the middle", input: "005f002k00.0i0k0", encoded: "785f002k000i0k", decoded: "5f002k00.0i0k"},

		{name: "one", input: "1", encoded: "711", decoded: "1"},
		{name: "ugly one", input: "+00001", encoded: "711", decoded: "1"},
		{name: "negative one", input: "-1", encoded: "3yy~", decoded: "-1"},
		{name: "ugly negative one", input: "-000001", encoded: "3yy~", decoded: "-1"},
		{name: "ugly positive int", input: "+00000123000", encoded: "76123", decoded: "123000"},
		{name: "ugly negative int", input: "-00000123000", encoded: "3tyxw~", decoded: "-123000"},
		{name: "fractional", input: "54321.12345", encoded: "755432112345", decoded: "54321.12345"},
		{name: "negative fractional", input: "-54321.12345", encoded: "3uuvwxyyxwvu~", decoded: "-54321.12345"},
		{name: "ugly fractional", input: "+00054321000.00012345000", encoded: "785432100000012345", decoded: "54321000.00012345"},
		{name: "ugly negative fractional", input: "-00054321000.00012345000", encoded: "3ruvwxyzzzzzzyxwvu~", decoded: "-54321000.00012345"},
		{name: "cowboy hat", input: "cowboy.hat", encoded: "76cowboyhat", decoded: "cowboy.hat"},
		{name: "negative cowboy hat", input: "-cowboy.hat", encoded: "3tnb3ob1ip6~", decoded: "-cowboy.hat"},
		{name: "maximum int length", input: "12345678901234567890123456789012345.1", encoded: "7z1123456789012345678901234567890123451", decoded: "12345678901234567890123456789012345.1"},
		{name: "maximum negative int length", input: "-12345678901234567890123456789012345.1", encoded: "30yyxwvutsrqzyxwvutsrqzyxwvutsrqzyxwvuy~", decoded: "-12345678901234567890123456789012345.1"},
		{name: "maximum fracleading zero count", input: "0.000000000000000000000000000000000004325430", encoded: "60y432543", decoded: "0.00000000000000000000000000000000000432543"},
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

func TestSortedness(t *testing.T) {
	step := 0.01
	prev := LessThanAny
	c := NewCodec()
	for i := -111111.0; i <= 111111.0; i++ {
		str := fmt.Sprintf("%3f", i*step)
		encoded, ok := c.Encode(str)
		if !ok {
			t.Fatal("Encoding failed for", i)
		}
		if prev >= encoded {
			t.Fatal("at", i*step, " ", prev, "is not smaller than", encoded)
		}
		prev = encoded
	}
}

func BenchmarkEncoding(b *testing.B) {
	step := 0.01
	c := NewCodec()
	to := float64(b.N / 2)
	from := -1 * to
	for i := from; i <= to; i++ {
		str := fmt.Sprintf("%3f", i*step)
		encoded, ok := c.Encode(str)
		if !ok {
			b.Fatal("Encoding failed for", i)
		}
		_, ok = c.Decode(encoded)
		if !ok {
			b.Fatal("Decoding failed for", encoded, "in", i)
		}
	}
}
