package conust

import (
	"strconv"
)

type base36Decoder struct{}

// NewBase36Decoder returns a Decoder that can decode base(36) Conust strings
func NewBase36Decoder() Decoder {
	return base36Decoder{}
}

// ToString decodes base(36) Conust string into a decimal number string
func (d base36Decoder) ToString(s string) (out string, ok bool) {
	// TODO
	return "", false
}

// ToInt32 decodes string into int32
// Successfulness of the decoding is signalled by the second return value. A failure is possible when the encoded number is out of the range of the int32 type.
func (d base36Decoder) ToInt32(s string) (i int32, ok bool) {
	if s == zeroOutput {
		return 0, true
	}
	intPart, _ := decodeStrings(s, true, flipDigit36)
	result, err := strconv.ParseInt(intPart, 36, 32)
	return int32(result), (err == nil)
}

// ToInt64 decodes string into int32
// Successfulness of the decoding is signalled by the second return value. A failure is possible when the encoded number is out of the range of the int32 type.
func (d base36Decoder) ToInt64(s string) (i int64, ok bool) {
	if s == zeroOutput {
		return 0, true
	}
	intPart, _ := decodeStrings(s, true, flipDigit36)
	result, err := strconv.ParseInt(intPart, 36, 64)
	return result, (err == nil)
}

// ToFloat32 decodes base(36) Conust string into a float 32
func (d base36Decoder) ToFloat32(s string) (f float32, ok bool) {
	// TODO
	return 0.0, false
}

// ToFloat64 decodes base(36) Conust string into a float 64
func (d base36Decoder) ToFloat64(s string) (f float64, ok bool) {
	// TODO
	return 0.0, false
}
