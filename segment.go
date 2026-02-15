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

const (
	segmentTypeText segmentTypeValue = iota
	segmentTypeNumber
)

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

func (s *segmentScanner) Next() (string, segmentTypeValue) {
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
		return s.input[startPos:digitPos], segmentTypeText
	}

	textEndPos := digitPos
	minusSignFound := false

	// adjust segment end index according to sign detection setting
	if digitPos > startPos {
		indexBeforeDigit := digitPos - 1
		byteBeforeDigit := s.input[indexBeforeDigit]

		if (byteBeforeDigit == plusByte && s.usePlusSign) ||
			(byteBeforeDigit == minusByte && s.useMinusSign) {
			textEndPos = indexBeforeDigit
			minusSignFound = (byteBeforeDigit == minusByte)
		}
	}

	// there is some text before the number (and optional sign character)
	if textEndPos > startPos {
		s.headPos = textEndPos
		return s.input[startPos:textEndPos], segmentTypeText
	}

	// the segment is a number segment

	// TODO parse number

	return "", segmentTypeText
}

func (s *segmentScanner) findFirstDigit() {
}

func (s *segmentScanner) xbyteType(b byte) byteTypeValue {
	if s.digitTester.isDigit(b) {
		return byteTypeDigit
	}

	if b == minusByte {
		if s.useMinusSign {
			return byteTypeMinus
		}
		return byteTypeText
	}

	if b == plusByte {
		if s.usePlusSign {
			return byteTypePlus
		}
		return byteTypeText
	}

	if b == s.thousandSeparator {
		if s.useThousandSeparator {
			return byteTypeThousandSeparator
		}
		return byteTypeText
	}

	if b == s.decimalSeparator {
		if s.useDecimalSeparator {
			return byteTypeDecimalSeparator
		}
		return byteTypeText
	}

	return byteTypeText
}

/*
func (s *segmentScanner) segmentTypeAtPos(pos int, previousSegmentType segmentTypeValue) segmentTypeValue {
	inputLen := len(s.input)

	if inputLen <= pos {
		return segmentTypeText
	}

	firstByteType := s.byteType(s.input[pos])

	if firstByteType == byteTypeText {
		return segmentTypeText
	}

	if firstByteType == byteTypeDigit {
		return segmentTypeNumber
	}

	if previousSegmentType == segmentTypeText {
		if firstByteType == byteTypePlus || firstByteType == byteTypeMinus {
			pos += 1

			if inputLen <= pos {
				return segmentTypeText
			}

			secondByteType := s.byteType(s.input[pos])

			if secondByteType == byteTypeDigit {
				return segmentTypeNumber
			}

			return segmentTypeText
		}
	}

	// previousSegmentType == segmentTypeNumber
	if firstByteType == byteTypeThousandSeparator || firstByteType == byteTypeDecimalSeparator {
		pos += 1

		if inputLen <= pos {
			return segmentTypeText
		}

		secondByteType := s.byteType(s.input[pos])

		if secondByteType == byteTypeDigit {
			return segmentTypeNumber
		}

		return segmentTypeText
	}

	return segmentTypeText
}
*/
