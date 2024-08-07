package aegis

import (
	"fmt"
	"io"
	"os"

	"github.com/minio/sio"
)

func decryptFile(inputPath, outputPath, password string) error {
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

	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(inputFile, salt); err != nil {
		return fmt.Errorf("error reading salt: %v", err)
	}

	key := deriveKey(password, salt)

	decryptedReader, err := sio.DecryptReader(inputFile, sio.Config{Key: key})
	if err != nil {
		return fmt.Errorf("error creating decrypted reader: %v", err)
	}

	if _, err := io.Copy(outputFile, decryptedReader); err != nil {
		return fmt.Errorf("error decrypting file: %v", err)
	}

	return nil
}
