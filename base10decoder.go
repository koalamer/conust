package conust

import "strconv"

type base10Decoder struct{}

// NewBase10Decoder returns a decoder that can decode base(10) Conust strings
func NewBase10Decoder() Decoder {
	return base10Decoder{}
}

// ToString decodes base(10) Conust string into a decimal number string
func (d base10Decoder) ToString(s string) (out string, ok bool) {
	//TODO
	return "", false
}

// ToInt32 decodes string into int32
// Successfulness of the decoding is signalled by the second return value. A failure is possible when the encoded number is out of the range of the int32 type.
func (d base10Decoder) ToInt32(s string) (i int32, ok bool) {
	if s == zeroOutput {
		return 0, true
	}
	intPart, _ := decodeStrings(s, true, flipDigit10)
	result, err := strconv.ParseInt(intPart, 10, 32)
	return int32(result), (err == nil)
}

// ToInt64 decodes string into int64
// Successfulness of the decoding is signalled by the second return value. A failure is possible when the encoded number is out of the range of the int64 type.
func (d base10Decoder) ToInt64(s string) (i int64, ok bool) {
	if s == zeroOutput {
		return 0, true
	}
	intPart, _ := decodeStrings(s, true, flipDigit10)
	result, err := strconv.ParseInt(intPart, 10, 64)
	return result, (err == nil)
}

// ToFloat32 decodes base10 Conust string into a float 32
func (d base10Decoder) ToFloat32(s string) (f float32, ok bool) {
	// TODO
	return 0.0, false
}

// ToFloat64 decodes base10 Conust string into a float 64
func (d base10Decoder) ToFloat64(s string) (f float64, ok bool) {
	// TODO
	return 0.0, false
}
