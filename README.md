# serial2net

![Go Logo](https://raw.githubusercontent.com/golang/go/master/doc/gopher/gopher.png)

A simple, single-binary Go application that acts as a relay between a serial port and a TCP port. It listens for a single incoming TCP connection and forwards any received data to the specified serial port.

## 🚀 Features

* **Serial to TCP Relay**: Forwards data received on a TCP connection directly to a serial port.
* **Single-Connection Mode**: Designed to handle only one TCP client at a time, ensuring dedicated access to the serial port.
* **Configurable**: All serial port and TCP settings are configurable via command-line flags, making it flexible for different environments.
* **Cross-Platform**: The Go compiler allows this tool to be easily compiled and run on Windows, Linux, and macOS.

## 📦 Getting Started

### Prerequisites

* **Go** (version 1.22 or higher)
* A serial device connected to your computer.

### Installation

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/your-username/your-repository.git](https://github.com/your-username/your-repository.git)
    cd your-repository
    ```

2.  **Install the necessary Go modules:**
    ```bash
    go mod tidy
    ```

3.  **Build the executable:**
    ```bash
    # For your current system
    go build -o relay

    # For Windows
    GOOS=windows go build -o relay.exe

    # For Linux
    GOOS=linux go build -o relay-linux
    ```

## ⚙️ Usage

The program is designed to be run from the command line with configurable flags.

### Command-line Flags

| Flag       | Type   | Default | Description                                 |
| :--------- | :----- | :------ | :------------------------------------------ |
| `--serial` | string | `COM3`  | The name of the serial port (e.g., `COM3`, `/dev/ttyS0`). |
| `--baud`   | int    | `9600`  | The baud rate for the serial connection.       |
| `--tcp`    | string | `:8000` | The TCP port to listen on (e.g., `:8000`).    |

### Example

To run the relay, connecting to a serial device on `COM4` with a baud rate of `115200` and listening for TCP connections on port `9000`:

```bash
./relay --serial COM4 --baud 115200 --tcp :9000
