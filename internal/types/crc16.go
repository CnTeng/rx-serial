package types

import "github.com/sigurn/crc16"

var crc16ARCTable = crc16.MakeTable(crc16.CRC16_ARC)

type CRC16Type struct {
	Scoop []Type
	IntType
}

func (c *CRC16Type) Serialize() []byte {
	b := []byte{}
	for _, t := range c.Scoop {
		b = append(b, t.Serialize()...)
	}
	c.Value = int(crc16.Checksum(b, crc16ARCTable))
	return c.IntType.Serialize()
}
