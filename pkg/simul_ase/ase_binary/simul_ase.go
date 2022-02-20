package ase_binary

/**
 * DWORD     File size
 * WORD      Magic number (0xA5E0)
 * WORD      Frames
 * WORD      Width in pixels
 * WORD      Height in pixels
 * WORD      Color depth (bits per pixel)
 *           32 bpp = RGBA
 *           16 bpp = Grayscale
 *           8 bpp = Indexed
 * DWORD     Flags:
 *           1 = Layer opacity has valid value
 * WORD      Speed (milliseconds between frame, like in FLC files)
 *           DEPRECATED: You should use the frame duration field
 *           from each frame header
 * DWORD     Set be 0
 * DWORD     Set be 0
 * BYTE      Palette entry (index) which represent transparent color
 *           in all non-background layers (only for Indexed sprites).
 * BYTE[3]   Ignore these bytes
 * WORD      Number of colors (0 means 256 for old sprites)
 * BYTE      Pixel width (pixel ratio is "pixel width/pixel height").
 *           If this or pixel height field is zero, pixel ratio is 1:1
 * BYTE      Pixel height
 * SHORT     X position of the grid
 * SHORT     Y position of the grid
 * WORD      Grid width (zero if there is no grid, grid size
 *           is 16x16 on Aseprite by default)
 * WORD      Grid height (zero if there is no grid)
 * BYTE[84]  For future (set to zero)
 */

type AsepriteHeader struct {
	FileSize       uint32
	MagicNumber    uint16 // A5E0
	Frames         uint16
	WidthInPixels  uint16
	HeightInPixels uint16
	ColorDepth     uint16
	Flags          uint32
	Speed          uint16 // deprecated, use frame duration from frame header
	// These are 0
	ignore1, ignore2 uint32
	PaletteEntry     byte // represents transparent color in all non-background layers (only for indexed sprites)
	// format defines these 3 bytes as ignored
	ignore3         [3]byte
	NumberOfColors  uint16
	PixelWidth      byte // (pixel ratio is "pxel width/pixel height")
	PixelHeight     byte
	XPositionOfGrid int16
	YPositionOfGrid int16
	GridWidth       uint16
	GridHeight      uint16
	// Future reserved
	reserved [84]byte
}

/**
 * DWORD     Bytes in this frame
 * WORD      Magic number (always 0xF1FA)
 * WORD      Old field which specifies the number of "chunks"
 *           in this frame. If this value is 0xFFFF, we might
 *           have more chunks to read in this frame
 *           (so we have to use the new field)
 * WORD      Frame duration (in milliseconds)
 * BYTE[2]   For future (set to zero)
 * DWORD     New field which specifies the number of "chunks"
 *           in this frame (if this is 0, use the old field)
 */

type AsepriteFrame struct {
	BytesThisFrame            uint32
	MagicNumber               uint16 // F1FA
	ChunksThisFrame           uint16 // If this value is FFFF there "may" be more chunks to read
	FrameDurationMilliseconds uint16
	reserved                  [2]byte
	ChunksThisFrameExt        uint32 // New field which specifies num chunks. If 0, use old field
}

/**
 * Then each chunk format is:
 * DWORD       Chunk size
 * WORD        Chunk type
 * BYTE[]      Chunk data
 */

type AsepriteFrameChunk struct {
	ChunkSize uint32
	ChunkType uint16
	ChunkData []byte
}

/**
 * Ignore this chunk if you find the new palette chunk (0x2019) Aseprite v1.1 saves both chunks 0x0004 and 0x2019 just for backward compatibility.
 * WORD        Number of packets
 * + For each packet
 *  BYTE      Number of palette entries to skip from the last packet (start from 0)
 *  BYTE      Number of colors in the packet (0 means 256)
 *  + For each color in the packet
 *    BYTE    Red (0-255)
 *    BYTE    Green (0-255)
 *    BYTE    Blue (0-255)
 */
type AsepriteOldPaletteChunk0004 struct {
	NumberOfPackets uint16
	// + for each packet
	Packets []
	//   + for each color in the packet

}

type AsepriteOldPaletteChunk0004Packet struct {
	NumPalletteEntriesToSkip byte // from thee last packet (start from 0)
	NumColorsInThePacket     byte // 0 means 256
}
