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
 * TILE: Tilemaps: Each tile can be a 8-bit (BYTE), 16-bit (WORD), or 32-bit (DWORD) value and there are masks
 * related to the meaning of each bit.
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

/**
 * Cel Chunk (0x2005)
 * This chunk determine where to put a cel in the specified layer/frame.
 *
 * WORD        Layer index (see NOTE.2)
 * SHORT       X position
 * SHORT       Y position
 * BYTE        Opacity level
 * WORD        Cel Type
 *           0 - Raw Image Data (unused, compressed image is preferred)
 *           1 - Linked Cel
 *           2 - Compressed Image
 *           3 - Compressed Tilemap
 * BYTE[7]     For future (set to zero)
 * + For cel type = 0 (Raw Image Data)
 *  WORD      Width in pixels
 *  WORD      Height in pixels
 *  PIXEL[]   Raw pixel data: row by row from top to bottom,
 *           for each scanline read pixels from left to right.
 * + For cel type = 1 (Linked Cel)
 *  WORD      Frame position to link with
 * + For cel type = 2 (Compressed Image)
 *  WORD      Width in pixels
 *  WORD      Height in pixels
 *  BYTE[]    "Raw Cel" data compressed with ZLIB method (see NOTE.3)
 * + For cel type = 3 (Compressed Tilemap)
 *  WORD      Width in number of tiles
 *  WORD      Height in number of tiles
 *  WORD      Bits per tile (at the moment it's always 32-bit per tile)
 *  DWORD     Bitmask for tile ID (e.g. 0x1fffffff for 32-bit tiles)
 *  DWORD     Bitmask for X flip
 *  DWORD     Bitmask for Y flip
 *  DWORD     Bitmask for 90CW rotation
 *  BYTE[10]  Reserved
 *  TILE[]    Row by row, from top to bottom tile by tile
 *            compressed with ZLIB method (see NOTE.3)
 */

type AsepriteCelChunk2005 struct {
	LayerIndex   uint16
	X, Y         int16
	OpacityLevel byte
	CelType      uint16
	future       [7]byte
	// + For cel type = 0 (Raw Image Data)
	WidthInPix, HeightInPix uint16
	RawPixData              []byte
	// + For cel type = 1 (Linked Cel)
	FramePosToLinkWith uint16
	// + For cel type = 2 (Compressed Image)
	// WORD      Width in pixels
	// WORD      Height in pixels
	RawCelData []byte // "Raw Cel" data compressed with ZLIB method (see NOTE.3)
	// + For cel type = 3 (Compressed Tilemap)
	WidthInTiles, HeightInTiles uint16
	BitsPerTile                 uint16 // (at the moment it's always 32-bit per tile)
	BitMaskForTileID            uint32 // (e.g. 0x1fffffff for 32-bit tiles)
	BitMaskForXFlip             uint32
	BitMaskForYFlip             uint32
	BitMaskFor90CWRot           uint32
	reserved                    [10]byte
	Tiles                       []byte // zlib data (see NOTE.3)
}

/**
 * Cel Extra Chunk (0x2006)
 * Adds extra information to the latest read cel.
 *
 * DWORD     Flags (set to zero)
 *            1 = Precise bounds are set
 * FIXED     Precise X position
 * FIXED     Precise Y position
 * FIXED     Width of the cel in the sprite (scaled in real-time)
 * FIXED     Height of the cel in the sprite
 * BYTE[16]  For future use (set to zero)
 */
type AsepriteCelExtraChunk2006 struct {
	Flags              uint32
	PreciseX, PreciseY uint32
	WidthCelInSprite   uint32
	HeightCelInSprite  uint32
	futureUse          [16]byte
}

/**
 * Color Profile Chunk (0x2007)
 * Color profile for RGB or grayscale values.
 *
 * WORD        Type
 *              0 - no color profile (as in old .aseprite files)
 *              1 - use sRGB
 *              2 - use the embedded ICC profile
 * WORD        Flags
 *              1 - use special fixed gamma
 * FIXED       Fixed gamma (1.0 = linear)
 *             Note: The gamma in sRGB is 2.2 in overall but it doesn't use
 *             this fixed gamma, because sRGB uses different gamma sections
 *             (linear and non-linear). If sRGB is specified with a fixed
 *             gamma = 1.0, it means that this is Linear sRGB.
 * BYTE[8]     Reserved (set to zero)
 * + If type = ICC:
 * DWORD     ICC profile data length
 * BYTE[]    ICC profile data. More info: http://www.color.org/ICC1V42.pdf
 */

type AsepriteColorProfileChunk2007 struct {
	Type       uint16
	Flags      uint16
	FixedGamma uint32
	reserved   [8]byte
	// + If type = ICC:
	ICCProfileDatLen uint32
	ICCProfileDat    []byte
}

/**
 * External Files Chunk (0x2008)
 * A list of external files linked with this file. It might be used to reference external palettes or tilesets.
 *
 * DWORD       Number of entries
 * BYTE[8]     Reserved (set to zero)
 * + For each entry
 *  DWORD     Entry ID (this ID is referenced by tilesets or palettes)
 *  BYTE[8]   Reserved (set to zero)
 *  STRING    External file name
 */

type AsepriteExternalFilesChunk2008 struct {
	NumEntries uint32
	reserved   [8]byte
	// + for each entry
	ExternalFile []AsepriteExternalFilesChunk2008Entry
}

type AsepriteExternalFilesChunk2008Entry struct {
	EntryID          uint32
	reserved         [8]byte
	ExternalFilename AsepriteString
}

/**
 * Mask Chunk (0x2016) DEPRECATED
 *
 * SHORT   X position
 * SHORT   Y position
 * WORD    Width
 * WORD    Height
 * BYTE[8] For future (set to zero)
 * STRING  Mask name
 * BYTE[]  Bit map data (size = height*((width+7)/8))
 *         Each byte contains 8 pixels (the leftmost pixels are
 *         packed into the high order bits)
 */

type AsepriteMaskChunk2016 struct {
	X, Y          int16
	Width, Height uint16
	future        [8]byte
	MaskName      AsepriteString
	BitMapData    []byte
}

/**
 * Path Chunk (0x2017)
 * Never used
 */

type AsepritePathChunk2017 struct{}

/**
 * Tags Chunk (0x2018)
 * After the tags chunk, you can write one user data chunk for each tag. E.g. if there are 10 tags, you can then write 10 user data chunks one for each tag.
 *
 * WORD        Number of tags
 * BYTE[8]     For future (set to zero)
 * + For each tag
 *  WORD      From frame
 *  WORD      To frame
 *  BYTE      Loop animation direction
 *              0 = Forward
 *              1 = Reverse
 *              2 = Ping-pong
 *  BYTE[8]   For future (set to zero)
 *  BYTE[3]   RGB values of the tag color
 *              Deprecated, used only for backward compatibility with Aseprite v1.2.x
 *              The color of the tag is the one in the user data field following
 *              the tags chunk
 *  BYTE      Extra byte (zero)
 *  STRING    Tag name
 */
type AsepriteTagsChunk2018 struct {
}

/**
 * Notes
 * NOTE.1
 * The child level is used to show the relationship of this layer with the last one read, for example:
 *
 * Layer name and hierarchy      Child Level
 * -----------------------------------------------
 * - Background                  0
 *   `- Layer1                   1
 * - Foreground                  0
 *   |- My set1                  1
 *   |  `- Layer2                2
 *   `- Layer3                   1
 *
 * NOTE.2
 * The layer index is a number to identify any layer in the sprite, for example:
 *
 * Layer name and hierarchy      Layer index
 * -----------------------------------------------
 * - Background                  0
 *   `- Layer1                   1
 * - Foreground                  2
 *   |- My set1                  3
 *   |  `- Layer2                4
 *   `- Layer3                   5
 *
 * NOTE.3
 * Details about the ZLIB and DEFLATE compression methods:
 * https://www.ietf.org/rfc/rfc1950
 * https://www.ietf.org/rfc/rfc1951
 * Some extra notes that might help you to decode the data: http://george.chiramattel.com/blog/2007/09/deflatestream-block-length-does-not-match.html
 */
