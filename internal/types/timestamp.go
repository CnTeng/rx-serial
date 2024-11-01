package types

import (
	"strings"
	"time"
)

type TimeStamp struct {
	Value time.Time
	*TypeMeta
}

func (c *TimeStamp) Parse(data []byte) error {
	c.Value = bcdToTime(data[:c.Len()])
	return nil
}

func (c *TimeStamp) Serialize() []byte {
	return timeToBCD(c.Value)
}

func (c *TimeStamp) Equal(other Type) bool {
	if oc, ok := other.(*TimeStamp); !ok {
		return false
	} else {
		return c.Value.Equal(oc.Value) && c.TypeMeta.Equal(oc.TypeMeta)
	}
}

func timeToBCD(t time.Time) []byte {
	formattedTime := t.Format("20060102150405")
	bcd := make([]byte, len(formattedTime)/2)
	for i := range bcd {
		bcd[i] = (formattedTime[i*2]-'0')<<4 | (formattedTime[i*2+1] - '0')
	}
	return bcd
}

func bcdToTime(bcd []byte) time.Time {
	var formattedTime strings.Builder
	formattedTime.Grow(len(bcd) * 2)
	for _, b := range bcd {
		formattedTime.WriteByte((b >> 4) + '0')
		formattedTime.WriteByte((b & 0x0F) + '0')
	}
	t, _ := time.Parse("20060102150405", formattedTime.String())
	return t
}
