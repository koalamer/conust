package conust

type byteTypeValue uint8

const (
	byteTypeText byteTypeValue = iota
	byteTypeDigit
	byteTypeMinus
	byteTypePlus
	byteTypeThousandSeparator
	byteTypeDecimalSeparator
)

type segmentTypeValue uint8

type segmentScanner struct {
	input string
	// config
	digitTester          digitTester
	thousandSeparator    byte
	decimalSeparator     byte
	useThousandSeparator bool
	useDecimalSeparator  bool
	usePlusSign          bool
	useMinusSign         bool
	// state
	headPos int
	// result info
	containsSingleNumber bool
}

func newSegmentScanner(
	radix int,
	thousandSeparator byte,
	decimalSeparator byte,
	useThousandSeparator bool,
	useDecimalSeparator bool,
	usePlusSign bool,
	useMinusSign bool,
) (*segmentScanner, error) {
	digitTester, err := newDigitTester(radix)
	if err != nil {
		return nil, err
	}

	return &segmentScanner{
		digitTester:          *digitTester,
		thousandSeparator:    thousandSeparator,
		decimalSeparator:     decimalSeparator,
		useThousandSeparator: useThousandSeparator,
		useDecimalSeparator:  useDecimalSeparator,
		usePlusSign:          usePlusSign,
		useMinusSign:         useMinusSign,
	}, nil
}

func (s *segmentScanner) Reset(input string) {
	s.input = input
	s.headPos = 0
	s.containsSingleNumber = true
}

func (s *segmentScanner) Next() (segment string, isNumber bool) {
	startPos := s.headPos
	digitPos := s.headPos
	inputLength := len(s.input)

	// find first digit
	for i := startPos; i < inputLength; i++ {
		if s.digitTester.isDigit(s.input[i]) {
			digitPos = i
			break
		}
	}

	// no digit was found
	if digitPos >= inputLength {
		s.headPos = digitPos
		return s.input[startPos:digitPos], false
	}

	numberSegmentStartPos := digitPos

	// adjust segment end index according to sign detection setting
	if digitPos > startPos {
		indexBeforeDigit := digitPos - 1
		byteBeforeDigit := s.input[indexBeforeDigit]

		if (byteBeforeDigit == plusByte && s.usePlusSign) ||
			(byteBeforeDigit == minusByte && s.useMinusSign) {
			numberSegmentStartPos = indexBeforeDigit
		}
	}

	// there is some text before the number (and optional sign character)
	if numberSegmentStartPos > startPos {
		s.headPos = numberSegmentStartPos
		return s.input[startPos:numberSegmentStartPos], false
	}

	// the segment is a number segment, find the end
	useDecSep := s.useDecimalSeparator
	useThSep := s.useThousandSeparator

	for i := digitPos + 1; i < inputLength; i++ {
		b := s.input[i]
		if s.digitTester.isDigit(b) {
			continue
		}

		if useThSep && b == s.thousandSeparator {
			afterSeparatorIndex := i + 1
			// separator is not followed by a digit, so the number ended before it
			if afterSeparatorIndex >= inputLength || !s.digitTester.isDigit(s.input[afterSeparatorIndex]) {
				s.headPos = i
				return s.input[numberSegmentStartPos:i], true
			}

			i = afterSeparatorIndex + 1
			continue
		}

		if useDecSep && b == s.decimalSeparator {
			afterSeparatorIndex := i + 1
			// separator is not followed by a digit, so the number ended before it
			if afterSeparatorIndex >= inputLength || !s.digitTester.isDigit(s.input[afterSeparatorIndex]) {
				s.headPos = i
				return s.input[numberSegmentStartPos:i], true
			}

			useThSep = false
			useDecSep = false

			i = afterSeparatorIndex + 1
			continue
		}

		return s.input[numberSegmentStartPos:i], true
	}

	return s.input[numberSegmentStartPos:], true
}

func (s *segmentScanner) findFirstDigit() {
}
