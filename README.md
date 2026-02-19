# Locked - ultra secure file locker
Locked is a CLI tool that allows you to lock your files into self-contained `.locker` archives. 

## Features
- **AES-256-GCM Encryption**: Files are encrypted using AES-256-GCM
- **Password Protection**:  Password-based key derivation using salt (argon2)
- **Self-contained**: `.locker` file contain all the metadata needed for decryption

## Installation
You can install Locked via the [releases page](https://github.com/PrathamGhaywat/locked/releases). There you will find the respective binaries for windows, linux and MacOS. It may break on MacOS, since I haven't tested it yet on that Platform.

### Building from source
Make sure you have Go installed on your system. Then, you can clone the repository and build the project using the following commands:
```bash
git clone https://github.com/PrathamGhaywat/locked.git
cd locked
go build -o locked . #windows users should use `go build -o locked.exe .`
```

## Usage
To lock a file:
```bash
./locked lock myfile.txt #creates a myfile.txt.locker file and prompts for a password
```
To unlock a file:
```bash
./locked unlock myfile.txt.locker #restores the original myfile.txt file and prompts for the password for decryption
```

To get an overview of all available commands and options:
```bash
./locked help
```
I recommend adding the `locked` binary to your system's PATH for easier access and also renaming the binary from e.g locked-windows-arm64.exe to just locked.exe or respective names for other platforms.