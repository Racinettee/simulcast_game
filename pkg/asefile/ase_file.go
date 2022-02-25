package ase_binary

import "io"

type AsepriteFile struct {
	Header AsepriteHeader
	Frames []AsepriteFrame
}

func (aseFile *AsepriteFile) Decode(r io.Reader) {
	aseFile.Header.Decode(r)

}
