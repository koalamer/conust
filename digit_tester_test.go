package conust

import (
	"slices"
	"strconv"
	"testing"
)

func TestisDigit36(t *testing.T) {
	var allBytes [256]byte

	for i := range allBytes {
		allBytes[i] = byte(i)
	}

	for radix := 2; radix <= 36; radix++ {
		testName := "radix " + strconv.Itoa(radix)

		t.Run(testName, func(t *testing.T) {
			lowercaseDigits := digits36[:radix]
			uppercaseDigits := uppercaseDigits36[:radix]

			tester, err := newDigitTester(radix)
			if err != nil {
				t.Fatal("instantiating digitTester returned an error")
			}
			if tester == nil {
				t.Fatal("instantiating digitTester returned nil tester")
			}

			for _, b := range allBytes {
				isLcDigit := slices.Contains(lowercaseDigits, b)
				isUcDigit := slices.Contains(uppercaseDigits, b)

				if (isLcDigit || isUcDigit) != tester.isDigit36(b) {
					t.Fatalf("wrong digit check for byte %v", b)
				}
			}
		})
	}
}
