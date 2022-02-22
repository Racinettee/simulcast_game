package ase_binary

// Some color definitions

// Just used here by certain packets
type AsepriteRGB24 struct {
	R, G, B byte
}

type AsepriteString struct {
	Length uint16
	Bytes  []byte
}

/**
 * ASE files use Intel (little-endian) byte order.
 *
 * BYTE: An 8-bit unsigned integer value
 * WORD: A 16-bit unsigned integer value
 * SHORT: A 16-bit signed integer value
 * DWORD: A 32-bit unsigned integer value
 * LONG: A 32-bit signed integer value
 * FIXED: A 32-bit fixed point (16.16) value
 * BYTE[n]: "n" bytes.
 * STRING:
 *   WORD: string length (number of bytes)
 *   BYTE[length]: characters (in UTF-8) The '\0' character is not included.
 * PIXEL: One pixel, depending on the image pixel format:
 *   RGBA: BYTE[4], each pixel have 4 bytes in this order Red, Green, Blue, Alpha.
 *   Grayscale: BYTE[2], each pixel have 2 bytes in the order Value, Alpha.
 *   Indexed: BYTE, each pixel uses 1 byte (the index).
 * TILE: Tilemaps: Each tile can be a 8-bit (BYTE), 16-bit (WORD), or 32-bit (DWORD) value and there are masks related to the meaning of each bit.
 *
 * File Header:
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
	Packets []AsepriteOldPaletteChunk0004Packet
}

type AsepriteOldPaletteChunk0004Packet struct {
	NumPalletteEntriesToSkip byte // from thee last packet (start from 0)
	NumColorsInThePacket     byte // 0 means 256
	//   + for each color in the packet
	Colors []AsepriteRGB24
}

/**
 * Old palette chunk (0x0011)
 * Ignore this chunk if you find the new palette chunk (0x2019)
 * WORD        Number of packets
 * + For each packet
 * BYTE      Number of palette entries to skip from the last packet (start from 0)
 * BYTE      Number of colors in the packet (0 means 256)
 * + For each color in the packet
 *   BYTE    Red (0-63)
 *   BYTE    Green (0-63)
 *   BYTE    Blue (0-63)
 */

type AsepritePaletteChunk0011 struct {
	NumberOfPackets uint16
	// + for each packet
	Packets []AsepriteOldPaletteChunk0011Packet
}

type AsepriteOldPaletteChunk0011Packet struct {
	NumPalletteEntriesToSkip byte // start from 0
	NumColorsInThePacket     byte // 0 means 256
	// + for each color in the packet
	Colors []AsepriteRGB24 // but using colors in the range 0-63
}

/**
 * Layer Chunk (0x2004)
 * In the first frame should be a set of layer chunks to determine the entire layers layout:
 * WORD  Flags:
 *        1 = Visible
 *        2 = Editable
 *        4 = Lock movement
 *        8 = Background
 *       16 = Prefer linked cels
 *       32 = The layer group should be displayed collapsed
 *       64 = The layer is a reference layer
 * WORD  Layer type
 *        0 = Normal (image) layer
 *        1 = Group
 *        2 = Tilemap
 * WORD  Layer child level (see NOTE.1)
 * WORD  Default layer width in pixels (ignored)
 * WORD  Default layer height in pixels (ignored)
 * WORD  Blend mode (always 0 for layer set)
 *        Normal         = 0
 *        Multiply       = 1
 *        Screen         = 2
 *        Overlay        = 3
 *        Darken         = 4
 *        Lighten        = 5
 *        Color Dodge    = 6
 *        Color Burn     = 7
 *        Hard Light     = 8
 *        Soft Light     = 9
 *        Difference     = 10
 *        Exclusion      = 11
 *        Hue            = 12
 *        Saturation     = 13
 *        Color          = 14
 *        Luminosity     = 15
 *        Addition       = 16
 *        Subtract       = 17
 *        Divide         = 18
 * BYTE  Opacity
 *        Note: valid only if file header flags field has bit 1 set
 * BYTE[3] For future (set to zero)
 * STRING  Layer name
 *  + If layer type = 2
 * DWORD   Tileset index
 */

type AsepriteLayerChunk2004 struct {
	Flags                uint16
	LayerType            uint16
	LayerChildLevel      uint16
	DefLayerWidthPixels  uint16 // (ignored)
	DefLayerHeightPixels uint16 // (ignored)
	BlendMode            uint16
	Opacity              byte // only valid if file headre flag field has bit 1 set
	forFuture            [3]byte
	LayerName            AsepriteString
	// + if layer type = 2
	TilesetIndex uint32
}
