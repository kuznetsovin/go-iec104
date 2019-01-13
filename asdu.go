package main

import (
	"encoding/binary"
)

const (
	M_ME_NB_1 = 11
)

func Parse_M_ME_NB(asdu []byte) (int, int) {
	//ЗАШЛУШКА ДЛЯ ПРИМЕРА

	//tp := int(asdu[0])
	//numix := int(asdu[1])
	//caustx := int(asdu[2])
	//ao := int(asdu[3])
	//addr := asdu[4:6]
	ioa := binary.LittleEndian.Uint16(asdu[6:9])
	value := binary.LittleEndian.Uint16(asdu[9:11])
	return int(ioa), int(value)
}
