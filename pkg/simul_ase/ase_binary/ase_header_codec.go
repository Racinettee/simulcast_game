package ase_binary

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"log"
)

var ble = binary.LittleEndian

func (aseHeader *AsepriteHeader) Decode(r io.Reader) {
	binary.Read(r, ble, &aseHeader.FileSize)
	binary.Read(r, ble, &aseHeader.MagicNumber)
	binary.Read(r, ble, &aseHeader.Frames)
	binary.Read(r, ble, &aseHeader.WidthInPixels)
	binary.Read(r, ble, &aseHeader.HeightInPixels)
	binary.Read(r, ble, &aseHeader.ColorDepth)
	binary.Read(r, ble, &aseHeader.Flags)
	binary.Read(r, ble, &aseHeader.Speed)
	binary.Read(r, ble, &aseHeader.ignore1)
	binary.Read(r, ble, &aseHeader.ignore2)
	binary.Read(r, ble, &aseHeader.PaletteEntry)
	binary.Read(r, ble, &aseHeader.ignore3)
	binary.Read(r, ble, &aseHeader.NumberOfColors)
	binary.Read(r, ble, &aseHeader.PixelWidth)
	binary.Read(r, ble, &aseHeader.PixelHeight)
	binary.Read(r, ble, &aseHeader.XPositionOfGrid)
	binary.Read(r, ble, &aseHeader.YPositionOfGrid)
	binary.Read(r, ble, &aseHeader.GridWidth)
	binary.Read(r, ble, &aseHeader.GridHeight)
	binary.Read(r, ble, &aseHeader.reserved)
}

func (aseHeader *AsepriteHeader) Encode(w io.Writer) {
	binary.Write(w, ble, &aseHeader.FileSize)
	binary.Write(w, ble, &aseHeader.MagicNumber)
	binary.Write(w, ble, &aseHeader.Frames)
	binary.Write(w, ble, &aseHeader.WidthInPixels)
	binary.Write(w, ble, &aseHeader.HeightInPixels)
	binary.Write(w, ble, &aseHeader.ColorDepth)
	binary.Write(w, ble, &aseHeader.Flags)
	binary.Write(w, ble, &aseHeader.Speed)
	binary.Write(w, ble, &aseHeader.ignore1)
	binary.Write(w, ble, &aseHeader.ignore2)
	binary.Write(w, ble, &aseHeader.PaletteEntry)
	binary.Write(w, ble, &aseHeader.ignore3)
	binary.Write(w, ble, &aseHeader.NumberOfColors)
	binary.Write(w, ble, &aseHeader.PixelWidth)
	binary.Write(w, ble, &aseHeader.PixelHeight)
	binary.Write(w, ble, &aseHeader.XPositionOfGrid)
	binary.Write(w, ble, &aseHeader.YPositionOfGrid)
	binary.Write(w, ble, &aseHeader.GridWidth)
	binary.Write(w, ble, &aseHeader.GridHeight)
	binary.Write(w, ble, &aseHeader.reserved)
}

func (aseFrame *AsepriteFrame) Decode(r io.Reader) {
	binary.Read(r, ble, &aseFrame.BytesThisFrame)
	binary.Read(r, ble, &aseFrame.MagicNumber)
	binary.Read(r, ble, &aseFrame.ChunksThisFrame)
	binary.Read(r, ble, &aseFrame.FrameDurationMilliseconds)
	binary.Read(r, ble, &aseFrame.reserved)
	binary.Read(r, ble, &aseFrame.ChunksThisFrameExt)
	//
	// Load n-amount of chunks
}

func (aseFrame AsepriteFrame) Encode(w io.Writer) {
	binary.Write(w, ble, &aseFrame.BytesThisFrame)
	binary.Write(w, ble, &aseFrame.MagicNumber)
	binary.Write(w, ble, &aseFrame.ChunksThisFrame)
	binary.Write(w, ble, &aseFrame.FrameDurationMilliseconds)
	binary.Write(w, ble, &aseFrame.reserved)
	binary.Write(w, ble, &aseFrame.ChunksThisFrameExt)
	//
	// Write n-amount of chunks
}

func (asePaletteChunk *AsepriteOldPaletteChunk0004) Decode(r io.Reader) {
	binary.Read(r, ble, &asePaletteChunk.NumberOfPackets)
	asePaletteChunk.Packets = make([]AsepriteOldPaletteChunk0004Packet, asePaletteChunk.NumberOfPackets)
	for x := 0; x < int(asePaletteChunk.NumberOfPackets); x += 1 {
		binary.Read(r, ble, &asePaletteChunk.Packets[x].NumPalletteEntriesToSkip)
		binary.Read(r, ble, &asePaletteChunk.Packets[x].NumColorsInThePacket)
		asePaletteChunk.Packets[x].Colors = make([]AsepriteRGB24, asePaletteChunk.Packets[x].NumColorsInThePacket)
		for y := 0; y < int(asePaletteChunk.Packets[x].NumColorsInThePacket); y += 1 {
			binary.Read(r, ble, &asePaletteChunk.Packets[x].Colors[y].R)
			binary.Read(r, ble, &asePaletteChunk.Packets[x].Colors[y].G)
			binary.Read(r, ble, &asePaletteChunk.Packets[x].Colors[y].B)
		}
	}
}

func (asePaletteChunk AsepriteOldPaletteChunk0004) Encode(w io.Writer) {
	binary.Write(w, ble, &asePaletteChunk.NumberOfPackets)
	for x := 0; x < int(asePaletteChunk.NumberOfPackets); x += 1 {
		binary.Write(w, ble, &asePaletteChunk.Packets[x].NumPalletteEntriesToSkip)
		binary.Write(w, ble, &asePaletteChunk.Packets[x].NumColorsInThePacket)
		for y := 0; y < int(asePaletteChunk.Packets[x].NumColorsInThePacket); y += 1 {
			binary.Write(w, ble, &asePaletteChunk.Packets[x].Colors[y].R)
			binary.Write(w, ble, &asePaletteChunk.Packets[x].Colors[y].G)
			binary.Write(w, ble, &asePaletteChunk.Packets[x].Colors[y].B)
		}
	}
}

func (asePaletteChunk *AsepritePaletteChunk0011) Decode(r io.Reader) {
	binary.Read(r, ble, &asePaletteChunk.NumberOfPackets)
	for x := 0; x < int(asePaletteChunk.NumberOfPackets); x += 1 {
		binary.Read(r, ble, &asePaletteChunk.Packets[x].NumPalletteEntriesToSkip)
		binary.Read(r, ble, &asePaletteChunk.Packets[x].NumColorsInThePacket)
		asePaletteChunk.Packets[x].Colors = make([]AsepriteRGB24, asePaletteChunk.Packets[x].NumColorsInThePacket)
		for y := 0; y < int(asePaletteChunk.Packets[x].NumColorsInThePacket); y += 1 {
			binary.Read(r, ble, &asePaletteChunk.Packets[x].Colors[y].R)
			binary.Read(r, ble, &asePaletteChunk.Packets[x].Colors[y].G)
			binary.Read(r, ble, &asePaletteChunk.Packets[x].Colors[y].B)
		}
	}
}

func (asePaletteChunk AsepritePaletteChunk0011) Encode(w io.Writer) {
	binary.Write(w, ble, &asePaletteChunk.NumberOfPackets)
	for x := 0; x < int(asePaletteChunk.NumberOfPackets); x += 1 {
		binary.Write(w, ble, &asePaletteChunk.Packets[x].NumPalletteEntriesToSkip)
		binary.Write(w, ble, &asePaletteChunk.Packets[x].NumColorsInThePacket)
		for y := 0; y < int(asePaletteChunk.Packets[x].NumColorsInThePacket); y += 1 {
			binary.Write(w, ble, &asePaletteChunk.Packets[x].Colors[y].R)
			binary.Write(w, ble, &asePaletteChunk.Packets[x].Colors[y].G)
			binary.Write(w, ble, &asePaletteChunk.Packets[x].Colors[y].B)
		}
	}
}

func (aseLayerChunk *AsepriteLayerChunk2004) Decode(r io.Reader) {
	binary.Read(r, ble, &aseLayerChunk.Flags)
	binary.Read(r, ble, &aseLayerChunk.LayerType)
	binary.Read(r, ble, &aseLayerChunk.LayerChildLevel)
	binary.Read(r, ble, &aseLayerChunk.DefLayerWidthPixels)
	binary.Read(r, ble, &aseLayerChunk.DefLayerHeightPixels)
	binary.Read(r, ble, &aseLayerChunk.BlendMode)
	binary.Read(r, ble, &aseLayerChunk.Opacity)
	binary.Read(r, ble, &aseLayerChunk.forFuture)
	binary.Read(r, ble, &aseLayerChunk.LayerName.Length)
	aseLayerChunk.LayerName.Bytes = make([]byte, aseLayerChunk.LayerName.Length)
	binary.Read(r, ble, aseLayerChunk.LayerName.Bytes)
	if aseLayerChunk.LayerType == 2 {
		binary.Read(r, ble, &aseLayerChunk.TilesetIndex)
	}
}

func (aseLayerChunk AsepriteLayerChunk2004) Encode(w io.Writer) {
	binary.Write(w, ble, &aseLayerChunk.Flags)
	binary.Write(w, ble, &aseLayerChunk.LayerType)
	binary.Write(w, ble, &aseLayerChunk.LayerChildLevel)
	binary.Write(w, ble, &aseLayerChunk.DefLayerWidthPixels)
	binary.Write(w, ble, &aseLayerChunk.DefLayerHeightPixels)
	binary.Write(w, ble, &aseLayerChunk.BlendMode)
	binary.Write(w, ble, &aseLayerChunk.Opacity)
	binary.Write(w, ble, &aseLayerChunk.forFuture)
	binary.Write(w, ble, &aseLayerChunk.LayerName.Length)
	binary.Write(w, ble, &aseLayerChunk.LayerName.Bytes)
	if aseLayerChunk.LayerType == 2 {
		binary.Write(w, ble, &aseLayerChunk.TilesetIndex)
	}
}

func (aseCelChunk *AsepriteCelChunk2005) Decode(r io.Reader) {
	binary.Read(r, ble, &aseCelChunk.LayerIndex)
	binary.Read(r, ble, &aseCelChunk.X)
	binary.Read(r, ble, &aseCelChunk.Y)
	binary.Read(r, ble, &aseCelChunk.OpacityLevel)
	binary.Read(r, ble, &aseCelChunk.CelType)
	binary.Read(r, ble, &aseCelChunk.future)
	switch aseCelChunk.CelType {
	case 0:
		binary.Read(r, ble, &aseCelChunk.WidthInPix)
		binary.Read(r, ble, &aseCelChunk.HeightInPix)
		bytesToAlloc := int(aseCelChunk.WidthInPix) * int(aseCelChunk.HeightInPix)
		switch aseCelChunk.parentHeader.ColorDepth {
		case 32:
			bytesToAlloc *= 4
		case 16:
			bytesToAlloc *= 2
		case 8:
			break
		}
		aseCelChunk.RawPixData = make([]byte, bytesToAlloc)
		binary.Read(r, ble, &aseCelChunk.RawPixData)
	case 1:
		binary.Read(r, ble, &aseCelChunk.FramePosToLinkWith)
	case 2:
		binary.Read(r, ble, &aseCelChunk.WidthInPix)
		binary.Read(r, ble, &aseCelChunk.HeightInPix)
		zreader, err := zlib.NewReader(r)
		if err != nil {
			log.Println(err)
		}
		byteBuff := bytes.NewBuffer([]byte{})
		io.Copy(byteBuff, zreader)
		aseCelChunk.RawCelData = byteBuff.Bytes()
	case 3:
		binary.Read(r, ble, &aseCelChunk.WidthInTiles)
		binary.Read(r, ble, &aseCelChunk.HeightInTiles)
		binary.Read(r, ble, &aseCelChunk.BitsPerTile)
		binary.Read(r, ble, &aseCelChunk.BitMaskForTileID)
		binary.Read(r, ble, &aseCelChunk.BitMaskForXFlip)
		binary.Read(r, ble, &aseCelChunk.BitMaskForYFlip)
		binary.Read(r, ble, &aseCelChunk.BitMaskFor90CWRot)
		binary.Read(r, ble, &aseCelChunk.reserved)
		zreader, err := zlib.NewReader(r)
		if err != nil {
			log.Println(err)
		}
		byteBuff := bytes.NewBuffer([]byte{})
		io.Copy(byteBuff, zreader)
		aseCelChunk.Tiles = byteBuff.Bytes()
	}
}

func (aseCelChunk *AsepriteCelChunk2005) Encode(w io.Writer) {
	binary.Write(w, ble, &aseCelChunk.LayerIndex)
	binary.Write(w, ble, &aseCelChunk.X)
	binary.Write(w, ble, &aseCelChunk.Y)
	binary.Write(w, ble, &aseCelChunk.OpacityLevel)
	binary.Write(w, ble, &aseCelChunk.CelType)
	binary.Write(w, ble, &aseCelChunk.future)
	switch aseCelChunk.CelType {
	case 0:
		binary.Write(w, ble, &aseCelChunk.WidthInPix)
		binary.Write(w, ble, &aseCelChunk.HeightInPix)
		bytesToAlloc := int(aseCelChunk.WidthInPix) * int(aseCelChunk.HeightInPix)
		switch aseCelChunk.parentHeader.ColorDepth {
		case 32:
			bytesToAlloc *= 4
		case 16:
			bytesToAlloc *= 2
		case 8:
			break
		}
		binary.Write(w, ble, &aseCelChunk.RawPixData)
	case 1:
		binary.Write(w, ble, &aseCelChunk.FramePosToLinkWith)
	case 2:
		binary.Write(w, ble, &aseCelChunk.WidthInPix)
		binary.Write(w, ble, &aseCelChunk.HeightInPix)
		var byteBuff bytes.Buffer
		zwriter := zlib.NewWriter(&byteBuff)
		zwriter.Write(aseCelChunk.RawCelData)
		w.Write(byteBuff.Bytes())
	case 3:
		binary.Write(w, ble, &aseCelChunk.WidthInTiles)
		binary.Write(w, ble, &aseCelChunk.HeightInTiles)
		binary.Write(w, ble, &aseCelChunk.BitsPerTile)
		binary.Write(w, ble, &aseCelChunk.BitMaskForTileID)
		binary.Write(w, ble, &aseCelChunk.BitMaskForXFlip)
		binary.Write(w, ble, &aseCelChunk.BitMaskForYFlip)
		binary.Write(w, ble, &aseCelChunk.BitMaskFor90CWRot)
		binary.Write(w, ble, &aseCelChunk.reserved)
		var byteBuff bytes.Buffer
		zwriter := zlib.NewWriter(&byteBuff)
		zwriter.Write(aseCelChunk.Tiles)
		w.Write(byteBuff.Bytes())
	}
}
