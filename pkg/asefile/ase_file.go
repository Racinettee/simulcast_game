package ase_binary

import "io"

type AsepriteFile struct {
	Header AsepriteHeader
	Frames []AsepriteFrame
}

func (aseFile *AsepriteFile) Decode(r io.Reader) {
	aseFile.Header.Decode(r)

	aseFile.Frames = make([]AsepriteFrame, aseFile.Header.Frames)

	for _, frame := range aseFile.Frames {
		frame.Decode(r)
	}
}

func (aseFile *AsepriteFile) Encode(w io.Writer) {
	aseFile.Header.Encode(w)

	aseFile.Frames = make([]AsepriteFrame, aseFile.Header.Frames)

	for _, frame := range aseFile.Frames {
		frame.Encode(w)
	}
}
