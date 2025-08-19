package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"go.bug.st/serial"
)

func main() {
	// Define flags for serial port and TCP settings
	serialPortName := flag.String("serial", "COM3", "Serial port name (e.g., COM3, /dev/ttyUSB0)")
	baudRate := flag.Int("baud", 9600, "Serial port baud rate")
	tcpPort := flag.String("tcp", ":8000", "TCP port to listen on (e.g., :8000)")

	// Parse the flags
	flag.Parse()

	// Open the serial port
	mode := &serial.Mode{
		BaudRate: *baudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
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
		if err != nil {
			if err == io.EOF {
				log.Printf("TCP connection from %s closed.", conn.RemoteAddr())
			} else {
				log.Printf("Error reading from TCP connection %s: %v", conn.RemoteAddr(), err)
			}
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