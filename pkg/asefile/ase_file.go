package asefile

import (
	"io"
)

type AsepriteFile struct {
	Header AsepriteHeader
	Frames []AsepriteFrame
}

func (aseFile *AsepriteFile) Decode(r io.Reader) error {
	aseFile.Header.Decode(r)
	aseFile.Frames = make([]AsepriteFrame, aseFile.Header.Frames)
	for _, frame := range aseFile.Frames {
		frame.parentHeader = &aseFile.Header
		err := frame.Decode(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (aseFile *AsepriteFile) Encode(w io.Writer) {
	aseFile.Header.Encode(w)

	aseFile.Frames = make([]AsepriteFrame, aseFile.Header.Frames)

	for _, frame := range aseFile.Frames {
		frame.Encode(w)
	}
}
