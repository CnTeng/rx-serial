package shell

import (
	"bufio"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/CnTeng/rx-serial/internal/data"
	"github.com/CnTeng/rx-serial/internal/message"
	"go.bug.st/serial"
)

func sendFile(reader *bufio.Reader, port serial.Port, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file stat:", err)
		return
	}

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		fmt.Println("Error hashing file:", err)
		return
	}
	md5sum := hash.Sum(nil)

	fileSize := stat.Size()
	fmt.Printf("File: %s, Size: %dbytes,  MD5: %s\n", fileName, fileSize, hex.EncodeToString(md5sum))

	buf := make([]byte, 101)
	buf[0] = 0
	copy(buf[1:81], fileName)
	binary.BigEndian.PutUint32(buf[81:85], uint32(fileSize))
	copy(buf[85:], md5sum[:16])

	fileData, _ := os.ReadFile(fileName)

	messages := message.GenerateMessage(0xC3, buf, 4096)
	messages = append(messages, message.GenerateMessage(0xC3, fileData, 4096)...)

	totalFrames := len(messages)

	for i, msg := range messages {
		msg.TotalFrames = uint16(totalFrames)
		msg.CurrentFrame = uint16(i + 1)
		msg.RefreshCRC()
		_, err := port.Write(msg.MarshalBinary())
		if err != nil {
			fmt.Println("Error writing to port:", err)
			return
		}

		readPort(port)
		fmt.Println("Press Enter to continue...")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		if strings.Contains(input, "\n") {
			continue
		}
	}
}

func readPort(port serial.Port) {
	time.Sleep(time.Millisecond * 200)

	c, err := data.NewConfig("test/config.json", "test/structs")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	buf := make([]byte, 1024)
	n, err := port.Read(buf)
	if err != nil {
		fmt.Println("Error reading from port:", err)
		return
	}

	fmt.Printf("Received %d: \n", n)

	if n == 0 {
		return
	}

	for _, b := range buf[:n] {
		fmt.Printf("%02x ", b)
	}
	fmt.Print("\n")

	d, err := c.Parse(buf)
	if err != nil {
		fmt.Println("Error parsing message:", err)
		return
	}

	fmt.Println(d)
}

func RawShell(port serial.Port) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		inputs := strings.Split(input, " ")
		if len(inputs) == 0 {
			readPort(port)
			continue
		}

		if inputs[0] == "quit" || inputs[0] == "q" {
			break
		} else if inputs[0] == "send" || inputs[0] == "s" {
			sendFile(reader, port, inputs[1])
			continue
		}

		bytes := make([]byte, len(inputs))

		for i, hexString := range inputs {
			byteValue, err := hex.DecodeString(hexString)
			if err != nil {
				fmt.Println("Error decoding hex string:", err)
				return
			}
			bytes[i] = byteValue[0]
		}

		messages := message.GenerateMessage(bytes[0], bytes[1:], 4096)

		for _, msg := range messages {
			_, err := port.Write(msg.MarshalBinary())
			if err != nil {
				fmt.Println("Error writing to port:", err)
				return
			}
		}

		readPort(port)
	}
}
