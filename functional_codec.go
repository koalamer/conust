package conust

import "strings"

type functionalCodec strings.Builder

// NewFunctionalCodec creates a functional kind of codec
func NewFunctionalCodec() (out Codec) {
	out = &functionalCodec{}
	return
}

func (fe *functionalCodec) Encode(in string) (out string, ok bool) {
	return "", false
}

func (fe *functionalCodec) Decode(in string) (out string, ok bool) {
	return "", false
}

func (fe *functionalCodec) EncodeInText(in string) (out string, ok bool) {
	return "", false
}
