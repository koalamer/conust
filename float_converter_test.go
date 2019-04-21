package conust

import (
	"fmt"
	"math"
	"testing"
)

type testCase struct {
	base             int
	decimalPrecision int
	input            float64
	formatted        string
	parsed           float64
}

func TestConversion(t *testing.T) {
	testCases := []testCase{

		// positive
		{base: 2, decimalPrecision: 3, input: 38.1234, formatted: "100110.000111111", parsed: 38.123},
		{base: 10, decimalPrecision: 3, input: 38.1234, formatted: "38.123", parsed: 38.123},
		{base: 16, decimalPrecision: 3, input: 38.1234, formatted: "26.1f8", parsed: 38.123},
		{base: 21, decimalPrecision: 3, input: 38.1234, formatted: "1h.2c5", parsed: 38.123},
		{base: 32, decimalPrecision: 3, input: 38.1234, formatted: "16.3u", parsed: 38.123},
		{base: 36, decimalPrecision: 3, input: 38.1234, formatted: "12.4f", parsed: 38.123},
		// negative
		{base: 2, decimalPrecision: 3, input: -38.1234, formatted: "-100110.000111111", parsed: -38.123},
		{base: 10, decimalPrecision: 3, input: -38.1234, formatted: "-38.123", parsed: -38.123},
		{base: 16, decimalPrecision: 3, input: -38.1234, formatted: "-26.1f8", parsed: -38.123},
		{base: 21, decimalPrecision: 3, input: -38.1234, formatted: "-1h.2c5", parsed: -38.123},
		{base: 32, decimalPrecision: 3, input: -38.1234, formatted: "-16.3u", parsed: -38.123},
		{base: 36, decimalPrecision: 3, input: -38.1234, formatted: "-12.4f", parsed: -38.123},
		// trailing zeros
		{base: 36, decimalPrecision: 9, input: 41.41, formatted: "15.ercyk6", parsed: 41.41},
		{base: 36, decimalPrecision: 9, input: -41.41, formatted: "-15.ercyk6", parsed: -41.41},
		// rounding
		{base: 17, decimalPrecision: 5, input: 323.000004, formatted: "120", parsed: 323.0},
		{base: 17, decimalPrecision: 5, input: 323.000006, formatted: "120.0000e", parsed: 323.00001},
		{base: 17, decimalPrecision: 5, input: 38.012344, formatted: "24.039ab", parsed: 38.01234},
		{base: 17, decimalPrecision: 5, input: 38.012346, formatted: "24.039b8", parsed: 38.01235},
		{base: 20, decimalPrecision: 2, input: 41.056, formatted: "21.14", parsed: 41.06},
		{base: 20, decimalPrecision: 2, input: -41.054, formatted: "-21.1", parsed: -41.05},
		{base: 20, decimalPrecision: 2, input: 41.0555, formatted: "21.14", parsed: 41.06},
		{base: 8, decimalPrecision: 2, input: 99.994, formatted: "143.773", parsed: 99.99},
		{base: 8, decimalPrecision: 2, input: -99.99, formatted: "-143.773", parsed: -99.99},
		{base: 16, decimalPrecision: 1, input: 15.96, formatted: "10", parsed: 16.0},
		{base: 16, decimalPrecision: 1, input: -15.96, formatted: "-10", parsed: -16.0},
	}

	for _, tc := range testCases {
		testName := "Base%d_DecPrecision%d_Input%f"
		t.Run(fmt.Sprintf(testName, tc.base, tc.decimalPrecision, tc.input), func(t *testing.T) {
			fc := NewFloatConverter(tc.base).WithDecimals(tc.decimalPrecision)
			f := fc.FormatFloat(tc.input)
			if f != tc.formatted {
				t.Fatalf("Expected formatted %s, got %s", tc.formatted, f)
			}

			p, ok := fc.ParseFloat(f)
			if !ok {
				t.Fatalf("Parsing %v returned ok = false", p)
			}

			allowedDelta := fc.precision / 2.0
			if math.Abs(tc.parsed-p) >= allowedDelta {
				t.Fatalf("Allowed delta %f exceeded for '%s'(%d), Actual delta: %f, Expected parsed %f, Actually parsed: %f",
					allowedDelta, f, tc.base, tc.parsed-p, tc.parsed, p)
			}
		})
	}
}
