package asefile

import (
	"fmt"
	"io"
)

type AsepriteFile struct {
	Header AsepriteHeader
	Frames []AsepriteFrame
}

func (aseFile *AsepriteFile) Decode(r io.Reader) error {
	fmt.Println("Decoding sprites...")
	aseFile.Header.Decode(r)
	aseFile.Frames = make([]AsepriteFrame, aseFile.Header.Frames)
	for _, frame := range aseFile.Frames {
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