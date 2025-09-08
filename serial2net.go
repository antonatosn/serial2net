package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"slices"
	"time"

	"go.bug.st/serial"
)

var acceptedBaudRates = []int{110, 300, 600, 1200, 2400, 4800, 9600, 14400, 19200, 38400, 57600, 115200, 230400, 460800, 921600}

func flagValidation() {
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	if !slices.Contains(acceptedBaudRates, *baudRate) {
		log.Fatalf("Invalid baud rate: %d", *baudRate)
	}
	if !slices.Contains([]int{5, 6, 7, 8}, *dataBits) {
		log.Fatalf("Invalid data bits: %d", *dataBits)
	}
	if *parity != "none" && *parity != "odd" && *parity != "even" {
		log.Fatalf("Invalid parity: %s", *parity)
	}
	if *stopBits != 1 && *stopBits != 2 {
		log.Fatalf("Invalid stop bits: %d", *stopBits)
	}
}

var (
	flags          = flag.NewFlagSet("serial2net", flag.ExitOnError)
	serialPortName = flags.String("serial", "COM3", "Serial port name (e.g., COM3, /dev/ttyUSB0)")
	baudRate       = flags.Int("baud", 9600, "Serial port baud rate (110, 300, 600, 1200, 2400, 4800, 9600, 14400, 19200, 38400, 57600, 115200, 230400, 460800, 921600)")
	dataBits       = flags.Int("bits", 8, "Serial port data bits (5, 6, 7, 8)")
	parity         = flags.String("parity", "none", "Serial port parity (none, odd, even)")
	stopBits       = flags.Int("stop", 1, "Serial port stop bits (1, 2)")
	tcpPort        = flags.String("tcp", ":8000", "TCP port to listen on (e.g., :8000)")
)

func main() {
	flags.Parse(os.Args[1:])

	flagValidation()

	// Open the serial port
	var parityMode serial.Parity
	switch *parity {
	case "none":
		parityMode = serial.NoParity
	case "odd":
		parityMode = serial.OddParity
	case "even":
		parityMode = serial.EvenParity
	}

	mode := &serial.Mode{
		BaudRate: *baudRate,
		DataBits: *dataBits,
		Parity:   parityMode,
		StopBits: serial.StopBits(*stopBits),
	}
	port, err := serial.Open(*serialPortName, mode)
	if err != nil {
		log.Fatalf("Failed to open serial port %s: %v", *serialPortName, err)
	}
	defer port.Close()
	log.Printf("Successfully opened serial port %s at %d baud.", *serialPortName, *baudRate)

	// Start listening on the TCP port
	listener, err := net.Listen("tcp", *tcpPort)
	if err != nil {
		log.Fatalf("Failed to listen on TCP port %s: %v", *tcpPort, err)
	}
	defer listener.Close()
	log.Printf("Listening for a single TCP connection on %s...", *tcpPort)

	// Accept a single incoming TCP connection
	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("Error accepting TCP connection: %v", err)
	}

	log.Printf("Accepted the TCP connection from %s", conn.RemoteAddr())

	// Handle the single connection
	handleConnection(conn, port)
}

func handleConnection(conn net.Conn, serPort serial.Port) {
	defer conn.Close()

	tcpReadBuffer := make([]byte, 1024)

	for {
		n, err := conn.Read(tcpReadBuffer)
		if err == io.EOF {
			log.Printf("TCP connection from %s closed.", conn.RemoteAddr())
			return
		}
		if err != nil {
			log.Printf("Error reading from TCP connection %s: %v", conn.RemoteAddr(), err)
			return
		}

		if n > 0 {
			dataToSerial := tcpReadBuffer[:n]
			log.Printf("Received %d bytes from TCP. Writing to serial port: %q", n, dataToSerial)
			_, err := serPort.Write(dataToSerial)
			if err != nil {
				log.Printf("Error writing to serial port: %v", err)
				return
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}
