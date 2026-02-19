package vault

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/PrathamGhaywat/locked/pkg/crypto"
)

const (
	//magic bytes to identify .locker files
	LockerMagic = "LCKR"
	//format version (current)
	LockerVersion = 1
)

// contains metadata about locked file.
type LockerHeader struct {
	Magic            [4]byte // "LCKR"
	Version          uint16
	Salt             [16]byte
	OriginalFilename string
	OriginalFileSize int64
}

// encrypts a file and creates .locker file
func CreateLocker(inputPath string, outputPath string, password string) error {
	//validate input
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return fmt.Errorf("failed to stat input file: %w", err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("use LockFolder for directories, not CreateLocker")
	}

	//generate salt
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	//derive encryption key from password
	key := crypto.DeriveKey(password, salt)

	//open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	//output .locker file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	//write header
	header := LockerHeader{
		Magic:            [4]byte{byte(LockerMagic[0]), byte(LockerMagic[1]), byte(LockerMagic[2]), byte(LockerMagic[3])},
		Version:          LockerVersion,
		Salt:             [16]byte(salt),
		OriginalFilename: filepath.Base(inputPath),
		OriginalFileSize: fileInfo.Size(),
	}

	err = writeHeader(outputFile, header)
	if err != nil {
		os.Remove(outputPath)
		return fmt.Errorf("failed to write header: %w", err)
	}

	err = crypto.EncryptStream(key, inputFile, outputFile)
	if err != nil {
		os.Remove(outputPath)
		return fmt.Errorf("failed to encrypt file: %w", err)
	}

	return nil
}

// decrypts .locker file
func OpenLocker(lockerPath string, outputPath string, password string) error {
	//open locker file
	lockerFile, err := os.Open(lockerPath)
	if err != nil {
		return fmt.Errorf("failed to open locker file: %w", err)
	}
	defer lockerFile.Close()

	//read + validate header
	header, err := readHeader(lockerFile)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	//derive key using stored salt
	key := crypto.DeriveKey(password, header.Salt[:])

	//create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	//decrypt and write contents
	err = crypto.DecryptStream(key, lockerFile, outputFile)
	if err != nil {
		os.Remove(outputPath)
		return fmt.Errorf("decryption failed: %w", err)
	}

	return nil
}

// write the LockerHeader to file
func writeHeader(w io.Writer, header LockerHeader) error {
	buffer := new(bytes.Buffer)

	//write magic
	_, err := buffer.Write(header.Magic[:])
	if err != nil {
		return fmt.Errorf("failed to write magic: %w", err)
	}

	//write version
	err = binary.Write(buffer, binary.LittleEndian, header.Version)
	if err != nil {
		return fmt.Errorf("failed to write version: %w", err)
	}

	//write salt
	_, err = buffer.Write(header.Salt[:])
	if err != nil {
		return fmt.Errorf("failed to write salt: %w", err)
	}

	//write filename length and filename
	err = binary.Write(buffer, binary.LittleEndian, uint32(len(header.OriginalFilename)))
	if err != nil {
		return fmt.Errorf("failed to write filename length: %w", err)
	}

	_, err = buffer.WriteString(header.OriginalFilename)
	if err != nil {
		return fmt.Errorf("failed to write filename: %w", err)
	}

	//write original file size
	err = binary.Write(buffer, binary.LittleEndian, header.OriginalFileSize)
	if err != nil {
		return fmt.Errorf("failed to write file size: %w", err)
	}

	//write buffer to writer
	_, err = w.Write(buffer.Bytes())
	return err
}

func readHeader(r io.Reader) (LockerHeader, error) {
	header := LockerHeader{}

	//read magic
	magic := make([]byte, 4)
	_, err := r.Read(magic)
	if err != nil {
		return header, fmt.Errorf("failed to read magic: %w", err)
	}
	copy(header.Magic[:], magic)

	// validate magic
	if string(header.Magic[:]) != LockerMagic {
		return header, fmt.Errorf("invalid locker file format")
	}

	//read version
	err = binary.Read(r, binary.LittleEndian, &header.Version)
	if err != nil {
		return header, fmt.Errorf("failed to read version: %w", err)
	}

	if header.Version != LockerVersion {
		return header, fmt.Errorf("unsupported locker version: %d", header.Version)
	}

	//read salt
	salt := make([]byte, 16)
	_, err = r.Read(salt)
	if err != nil {
		return header, fmt.Errorf("failed to read salt: %w", err)
	}
	copy(header.Salt[:], salt)

	// read filename length
	var filenameLen uint32
	err = binary.Read(r, binary.LittleEndian, &filenameLen)
	if err != nil {
		return header, fmt.Errorf("failed to read filename length: %w", err)
	}

	// read filename
	filenameBuf := make([]byte, filenameLen)
	_, err = r.Read(filenameBuf)
	if err != nil {
		return header, fmt.Errorf("failed to read filename: %w", err)
	}
	header.OriginalFilename = string(filenameBuf)

	// read file size
	err = binary.Read(r, binary.LittleEndian, &header.OriginalFileSize)
	if err != nil {
		return header, fmt.Errorf("failed to read file size: %w", err)
	}

	return header, nil
}

// GetOriginalFilename reads a .locker file and returns the original filename from its header.
func GetOriginalFilename(lockerPath string) (string, error) {
	lockerFile, err := os.Open(lockerPath)
	if err != nil {
		return "", fmt.Errorf("failed to open locker file: %w", err)
	}
	defer lockerFile.Close()

	header, err := readHeader(lockerFile)
	if err != nil {
		return "", fmt.Errorf("failed to read header: %w", err)
	}

	return header.OriginalFilename, nil
}
