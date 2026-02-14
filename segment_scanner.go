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
	input                string
	digitTester          digitTester
	nextSegmentStartPos  int
	headPos              int
	nextSegmentType      segmentTypeValue
	lastByteType         byteTypeValue
	thousandSeparator    byte
	decimalSeparator     byte
	useThousandSeparator bool
	useDecimalSeparator  bool
	usePlusSign          bool
	useMinusSign         bool
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
		input:                "",
		digitTester:          *digitTester,
		nextSegmentStartPos:  0,
		headPos:              0,
		nextSegmentType:      segmentTypeText,
		lastByteType:         byteTypeText,
		thousandSeparator:    thousandSeparator,
		decimalSeparator:     decimalSeparator,
		useThousandSeparator: useThousandSeparator,
		useDecimalSeparator:  useDecimalSeparator,
		usePlusSign:          usePlusSign,
		useMinusSign:         useMinusSign,
		containsSingleNumber: false,
	}, nil
}

func (s *segmentScanner) Reset(input string) {
	s.input = input
	s.nextSegmentStartPos = 0
	s.headPos = 0
	s.nextSegmentType = s.segmentTypeAtPos(0, segmentTypeText)
	s.lastByteType = byteTypeText
	s.containsSingleNumber = false
}

func (s *segmentScanner) Next() (string, segmentTypeValue) {
	if s.nextSegmentStartPos >= len(s.input) {
		return "", segmentTypeText
	}

	// nextSegmentType is already set
	s.headPos = s.nextSegmentStartPos

	for i := s.headPos; i < len(s.input); i++ {
		switch s.byteType(s.input[i]) {
		case byteTypeDigit:
			if s.lastByteType == byteTypeDigit {
				continue
			}
		}
	}

	return "", segmentTypeText
}

func (s *segmentScanner) byteType(b byte) byteTypeValue {
	if s.digitTester.isDigit(b) {
		return byteTypeDigit
	}

	if isMinusByte(b) && s.useMinusSign {
		return byteTypeMinus
	}

	if isPlusByte(b) && s.usePlusSign {
		return byteTypePlus
	}

	if b == s.thousandSeparator && s.useThousandSeparator {
		return byteTypeThousandSeparator
	}

	if b == s.decimalSeparator && s.useDecimalSeparator {
		return byteTypeDecimalSeparator
	}

	return byteTypeText
}

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
