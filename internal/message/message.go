package message

import (
	"encoding/binary"
	"time"

	"github.com/sigurn/crc16"
)

const (
	STX = 0x02
	ETX = 0x03
	DLE = 0x10
)

var CRC16ARCTable = crc16.MakeTable(crc16.CRC16_ARC)

type MessageData struct {
	Cmd          byte
	RepeatNo     byte
	CmdDataLen   uint16
	CmdData      []byte
	ResponseCode byte
	RFU          [8]byte
}

type Message struct {
	DataLen      uint16
	TotalFrames  uint16
	CurrentFrame uint16
	Time         [7]byte
	Data         *MessageData
	CRC          uint16
}

func NewMessageData(cmd byte, cmdData []byte) *MessageData {
	return &MessageData{
		Cmd:        cmd,
		RepeatNo:   0,
		CmdDataLen: uint16(len(cmdData)),
		CmdData:    cmdData,
	}
}

func toBCD(n int) byte {
	return byte((n/10)<<4 | n%10)
}

func getBCDTime() [7]byte {
	now := time.Now()

	return [7]byte{
		toBCD(now.Year() / 100),
		toBCD(now.Year() % 100),
		toBCD(int(now.Month())),
		toBCD(now.Day()),
		toBCD(now.Hour()),
		toBCD(now.Minute()),
		toBCD(now.Second()),
	}
}

func (m *Message) RefreshCRC() {
	crcBuf := make([]byte, 6)
	binary.BigEndian.PutUint16(crcBuf[0:2], m.DataLen)
	binary.BigEndian.PutUint16(crcBuf[2:4], m.TotalFrames)
	binary.BigEndian.PutUint16(crcBuf[4:6], m.CurrentFrame)
	crcBuf = append(crcBuf, m.Time[:]...)
	data := m.Data.MarshalBinary()
	crcBuf = append(crcBuf, data...)

	m.CRC = crc16.Checksum(crcBuf, CRC16ARCTable)
}

func NewMessage(totalFrames uint16, currentFrame uint16, cmd byte, cmdData []byte) *Message {
	time := getBCDTime()
	messageData := NewMessageData(cmd, cmdData)
	dataLen := messageData.Len()

	crcBuf := make([]byte, 6)

	binary.BigEndian.PutUint16(crcBuf[0:2], dataLen)
	binary.BigEndian.PutUint16(crcBuf[2:4], totalFrames)
	binary.BigEndian.PutUint16(crcBuf[4:6], currentFrame)
	crcBuf = append(crcBuf, time[:]...)
	data := messageData.MarshalBinary()
	crcBuf = append(crcBuf, data...)

	crc := crc16.Checksum(crcBuf, CRC16ARCTable)

	return &Message{
		DataLen:      dataLen,
		TotalFrames:  totalFrames,
		CurrentFrame: currentFrame,
		Time:         getBCDTime(),
		Data:         NewMessageData(cmd, cmdData),
		CRC:          crc,
	}
}

func (md *MessageData) MarshalBinary() []byte {
	var buf []byte

	buf = append(buf, md.Cmd)
	buf = append(buf, md.RepeatNo)

	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, uint16(md.CmdDataLen))
	buf = append(buf, tmp...)

	buf = append(buf, md.CmdData...)
	buf = append(buf, md.ResponseCode)
	buf = append(buf, md.RFU[:]...)

	return buf
}

func (md *MessageData) Len() uint16 {
	return 1 + 1 + 2 + md.CmdDataLen + 1 + 8
}

func (m *Message) MarshalBinary() []byte {
	buf := make([]byte, 6)

	binary.BigEndian.PutUint16(buf[0:2], m.DataLen)
	binary.BigEndian.PutUint16(buf[2:4], m.TotalFrames)
	binary.BigEndian.PutUint16(buf[4:6], m.CurrentFrame)
	buf = append(buf, m.Time[:]...)
	data := m.Data.MarshalBinary()
	buf = append(buf, data...)

	tmp := make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, m.CRC)
	buf = append(buf, tmp...)

	var escapedBuf []byte
	for _, b := range buf {
		if b == STX || b == ETX || b == DLE {
			escapedBuf = append(escapedBuf, DLE)
		}
		escapedBuf = append(escapedBuf, b)
	}

	return append([]byte{STX}, append(escapedBuf, ETX)...)
}

func GenerateMessage(cmd byte, cmdData []byte, maxLen int) []*Message {
	messagesLen := (len(cmdData) + maxLen - 1) / maxLen
	messages := make([]*Message, 0, messagesLen)

	for i := 0; i < messagesLen; i++ {
		start := i * maxLen
		end := start + maxLen
		if end > len(cmdData) {
			end = len(cmdData)
		}
		messages = append(messages, NewMessage(uint16(messagesLen), uint16(i+1), cmd, cmdData[start:end]))
	}

	return messages
}
