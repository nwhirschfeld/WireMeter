# WireMeter

WireMeter is a powerful toolkit developed in Golang using the gopacket library for testing network devices to evaluate delays and packet loss.

## Features

- **Delay Testing**: Measure the latency of network devices to assess response times.
- **Packet Loss Detection**: Identify packet loss across network connections to diagnose potential issues.
- **User-Friendly Web Interface**: WireMeter features a responsive web interface, making it easy to conduct tests and interpret results from any modern web browser.
- **Detailed Reports**: Measurements can be exported as raw csv data and svg graphs.

## Requirements

WireMeter is designed to run on platforms with at least three native Ethernet interfaces, such as the now-discontinued PC Engines APUv2.

## Installation

The includes script `install-alpine.sh` can be used to install WireMeter on a system running Alpine Linux.
It can be used as reference for installation on other systems, too.

## Usage

WireMeter provides a web interface on port 3000.

### Command Line Options

WireMeter can be configured to match your setup using its command line options:

```
-r string
Receiving interface (default "enp1")
-s string
Sending interface (default "enp2")
-sleep int
Time to wait in between send requests (good values may depend on your network interface) (default 5)
```
Important: Please ensure that the specified sending (-s) and receiving (-r) interfaces are not utilized for any other purposes while WireMeter is running to avoid interference with the testing process.

## Contributing

Contributions to WireMeter are welcome! If you have suggestions for improvements, bug fixes, or new features, please feel free to submit a pull request. 
Before making significant changes, it's advisable to open an issue to discuss the proposed modifications.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.