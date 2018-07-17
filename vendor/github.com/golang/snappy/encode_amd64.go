







package snappy




func emitLiteral(dst, lit []byte) int




func emitCopy(dst []byte, offset, length int) int




func extendMatch(src []byte, i, j int) int




func encodeBlock(dst, src []byte) (d int)
