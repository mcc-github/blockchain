







package snappy 

import (
	"hash/crc32"
)


const (
	tagLiteral = 0x00
	tagCopy1   = 0x01
	tagCopy2   = 0x02
	tagCopy4   = 0x03
)

const (
	checksumSize    = 4
	chunkHeaderSize = 4
	magicChunk      = "\xff\x06\x00\x00" + magicBody
	magicBody       = "sNaPpY"

	
	
	
	
	
	
	
	
	maxBlockSize = 65536

	
	
	
	
	maxEncodedLenOfMaxBlockSize = 76490

	obufHeaderLen = len(magicChunk) + checksumSize + chunkHeaderSize
	obufLen       = obufHeaderLen + maxEncodedLenOfMaxBlockSize
)

const (
	chunkTypeCompressedData   = 0x00
	chunkTypeUncompressedData = 0x01
	chunkTypePadding          = 0xfe
	chunkTypeStreamIdentifier = 0xff
)

var crcTable = crc32.MakeTable(crc32.Castagnoli)



func crc(b []byte) uint32 {
	c := crc32.Update(0, crcTable, b)
	return uint32(c>>15|c<<17) + 0xa282ead8
}
