package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CnTeng/rx-serial/internal/shell"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

var (
	baudrate int
	databits int
	stopbits int
	parity   string
	portName string
)

var rootCmd = &cobra.Command{
	Use:   "rx-serial",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var stopBits serial.StopBits
		switch stopbits {
		case 1:
			stopBits = serial.OneStopBit
		case 2:
			stopBits = serial.TwoStopBits
		default:
			log.Fatalf("Invalid stopbits value: %d", stopbits)
		}

		var parityMode serial.Parity
		switch parity {
		case "none":
			parityMode = serial.NoParity
		case "odd":
			parityMode = serial.OddParity
		case "even":
			parityMode = serial.EvenParity
		default:
			log.Fatalf("Invalid parity value: %s", parity)
		}

		mode := &serial.Mode{
			BaudRate: baudrate,
			DataBits: databits,
			Parity:   parityMode,
			StopBits: stopBits,
		}

		port, err := serial.Open(portName, mode)
		if err != nil {
			log.Fatal(err)
		}
		defer port.Close()
		fmt.Printf("Port %s opened with baudrate %d, databits %d, stopbits %d, and parity %s\n", portName, baudrate, databits, stopbits, parity)

		err = port.SetDTR(false)
		if err != nil {
			log.Fatal(err)
		}

		err = port.SetReadTimeout(time.Millisecond * 100)
		if err != nil {
			log.Fatal(err)
		}

		shell.RawShell(port)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVar(&baudrate, "baudrate", 115200, "baudrate")
	rootCmd.Flags().IntVar(&databits, "databits", 8, "databits")
	rootCmd.Flags().IntVar(&stopbits, "stopbits", 1, "stopbits (1 or 2)")
	rootCmd.Flags().StringVar(&parity, "parity", "none", "parity (none, odd, even)")
	rootCmd.Flags().StringVar(&portName, "port", "/dev/ttyUSB1", "serial port name")
}
