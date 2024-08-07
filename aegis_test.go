package aegis

import (
	"crypto/rand"
	"fmt"
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

func BenchmarkEncrypt(b *testing.B) {
	sizes := []int64{1024, 1024 * 1024, 10 * 1024 * 1024, 100 * 1024 * 1024, 1024 * 1024 * 1024}
	password := "benchmarkpassword"

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dB", size), func(b *testing.B) {
			tmpFile, err := createTempFile(size)
			if err != nil {
				b.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			outputFile := tmpFile.Name() + ".enc"
			defer os.Remove(outputFile)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := encryptFile(tmpFile.Name(), outputFile, password); err != nil {
					b.Fatalf("encryption failed: %v", err)
				}
			}
			b.StopTimer()
		})
	}
}

func BenchmarkDecrypt(b *testing.B) {
	sizes := []int64{1024, 1024 * 1024, 10 * 1024 * 1024, 100 * 1024 * 1024, 1024 * 1024 * 1024}
	password := "benchmarkpassword"

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%dB", size), func(b *testing.B) {
			tmpFile, err := createTempFile(size)
			if err != nil {
				b.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			encFile := tmpFile.Name() + ".enc"
			if err := encryptFile(tmpFile.Name(), encFile, password); err != nil {
				b.Fatalf("encryption failed: %v", err)
			}
			defer os.Remove(encFile)

			decFile := tmpFile.Name() + ".dec"
			defer os.Remove(decFile)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := decryptFile(encFile, decFile, password); err != nil {
					b.Fatalf("decryption failed: %v", err)
				}

				originalData, err := os.ReadFile(tmpFile.Name())
				if err != nil {
					b.Fatalf("failed to read original file: %v", err)
				}

				decryptedData, err := os.ReadFile(decFile)
				if err != nil {
					b.Fatalf("failed to read decrypted file: %v", err)
				}

				if !compareFiles(originalData, decryptedData) {
					b.Fatalf("decrypted data does not match original data")
				}
			}
			b.StopTimer()
		})
	}
}

func compareFiles(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
