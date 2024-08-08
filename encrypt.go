package aegis

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"github.com/minio/sio"
)

func EncryptFile(inputPath, outputPath, password string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	salt := make([]byte, SALTSIZE)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("error generating salt: %v", err)
	}

	key := DeriveKey(password, salt)

	if _, err := outputFile.Write(salt); err != nil {
		return fmt.Errorf("error writing salt: %v", err)
	}

	encryptedWriter, err := sio.EncryptWriter(outputFile, sio.Config{Key: key})
	if err != nil {
		return fmt.Errorf("error creating encrypted writer: %v", err)
	}

	if _, err := io.Copy(encryptedWriter, inputFile); err != nil {
		return fmt.Errorf("error encrypting file: %v", err)
	}

	if err := encryptedWriter.Close(); err != nil {
		return fmt.Errorf("error closing encrypted writer: %v", err)
	}

	return nil
}
