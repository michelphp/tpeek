# TPeek 🚀

A lightweight, zero-dependency TCP transparent proxy designed for deep traffic inspection and debugging. Built with Go for high performance and low resource footprint.

## Features

* **Bidirectional Inspection**: Real-time visualization of Client and Server data streams.
* **Visual Contrast**: Distinctive color-coded frames (🔵 Client / 🟢 Server) for instant source identification.
* **Safe for Production**: Forwards untouched binary data while sanitizing terminal output (prevents terminal crashes from binary noise).
* **Hex & Plain Text**: Toggle between raw hex dumps and readable text modes.
* **Zero Dependencies**: Single binary, no runtime required (unlike Python or Node.js tools).

## Installation

### From Source

```bash
sudo wget https://github.com/michelphp/tpeek/releases/download/1.0.0/tpeek-linux -O /usr/local/bin/tpeek
sudo chmod +x /usr/local/bin/tpeek
```


## Usage

### Basic Debugging (e.g., MySQL)

```bash
tpeek -l 0.0.0.0:8000 -t 127.0.0.1:3306

```

Point your application to `localhost:8000` to inspect the traffic flowing to your database.

### Hexadecimal Mode (Binary Protocols)

```bash
tpeek -l 0.0.0.0:8000 -t 127.0.0.1:6379 -hex

```

## Command Line Flags

| Flag | Description | Default |
| --- | --- | --- |
| `-l` | Local address and port to listen on |  |
| `-t` | Target service address and port |  |
| `-hex` | Enables full hexadecimal dump mode | `false` |


## Why TPeek?

TPeek is designed to observe the dialogue between two applications without the complexity of a network sniffer.

* **Flow Identification**: Each data block is prefixed by its source (Client or Server). This allows you to follow the "request-response" logic directly in the console.
* **Fragmentation Awareness**: TPeek displays data block by block, exactly as received by the `Read()` syscall. This shows how the application segments its messages instead of displaying a continuous, unstructured text stream.
* **Millisecond Timestamps**: Every read operation is marked with a precise timestamp. This makes it possible to measure the processing delay between an outgoing request and the incoming response.
* **Application-Level View**: The program displays the final data as reassembled by the operating system. You see exactly what the software receives in its buffer, without the overhead of lower-level network layers.


## License

GNU Affero General Public License v3.0 (AGPL-3.0)




