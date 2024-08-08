package aegis

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
)

func createTempFile(size int64) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "aegis_test_")
	if err != nil {
		return nil, err
	}

	data := make([]byte, size)
	if _, err = rand.Read(data); err != nil {
		return nil, err
	}

	if _, err = tmpFile.Write(data); err != nil {
		return nil, err
	}

	if err = tmpFile.Close(); err != nil {
		return nil, err
	}

	return tmpFile, nil
}

func BenchmarkEncryptDecrypt(b *testing.B) {
	sizes := []int64{1024, 1024 * 1024, 10 * 1024 * 1024, 100 * 1024 * 1024, 1024 * 1024 * 1024}
	password := "testpassword"
	for _, size := range sizes {
		tmpFile, err := createTempFile(size)
		if err != nil {
			b.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		outputFile := tmpFile.Name() + ".enc"
		defer os.Remove(outputFile)

		// Run encryption benchmark
		b.Run(fmt.Sprintf("%dB", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := EncryptFile(tmpFile.Name(), outputFile, password); err != nil {
					b.Fatalf("encryption failed: %v", err)
				}
			}
			b.StopTimer()
		})

		decFile := tmpFile.Name() + ".dec"
		defer os.Remove(decFile)
		// Run decryption benchmark
		b.Run(fmt.Sprintf("%dB", size), func(b *testing.B) {
			encFile := tmpFile.Name() + ".enc"
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := DecryptFile(encFile, decFile, password); err != nil {
					b.Fatalf("decryption failed: %v", err)
				}

				if !compareFiles(tmpFile.Name(), decFile) {
					b.Fatalf("decrypted data does not match original data")
				}
			}
			b.StopTimer()
		})
	}
}

const chunkSize = 64000

func compareFiles(file1, file2 string) bool {
	f1, err := os.Open(file1)
	if err != nil {
		log.Fatal(err)
	}
	f1i, err := f1.Stat()
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	f2, err := os.Open(file2)
	if err != nil {
		log.Fatal(err)
	}
	f2i, err := f2.Stat()
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	if f1i.Size() != f2i.Size() {
		return false
	}

	for {
		b1 := make([]byte, chunkSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}
