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
# if you are on linux, you need to make it executable: chmod +x locked
```

## Usage
### Windows
You have to rename the binary from e.g `locked-windows-arm64.exe`/`locked-windows-amd64.exe`/`locked-windows-386.exe` to `locked.exe` and then you can use it as following:

Locking:
```bash
locked.exe lock myfile.txt #creates a myfile.txt.locker file and prompts for a password
```

Unlocking:
```bash
locked.exe unlock myfile.txt.locker #restores the original myfile.txt file and prompts for the password for decryption
```

### Linux
To lock a file:
```bash
locked lock myfile.txt #creates a myfile.txt.locker file and prompts for a password
```
To unlock a file:
```bash
locked unlock myfile.txt.locker #restores the original myfile.txt file and prompts for the password for decryption
```

To get an overview of all available commands and options:
```bash
locked help
```
I recommend adding the `locked` binary (you have to give the folder path where the binary is located) to your system's PATH for easier access and also renaming the binary from e.g locked-windows-arm64.exe to just locked.exe or respective names for other platforms.