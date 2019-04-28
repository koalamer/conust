package conust

import "strings"

type sliceyCodec struct {
	builder strings.Builder
}

// NewSliceyCodec creates a slicey kind of codec
func NewSliceyCodec() (out Codec) {
	out = &sliceyCodec{}
	return
}

func (fe *sliceyCodec) Encode(in string) (out string, ok bool) {
	return "", false
}

func (fe *sliceyCodec) Decode(in string) (out string, ok bool) {
	return "", false
}

func (fe *sliceyCodec) EncodeInText(in string) (out string, ok bool) {
	return "", false
}
