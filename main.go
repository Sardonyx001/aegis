package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/minio/sio"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/pbkdf2"
)

// Salt size in bytes.
const saltSize = 16

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 4096, chacha20poly1305.KeySize, sha256.New)
}

func encryptFile(inputPath, outputPath, password string) error {
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
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("error generating salt: %v", err)
	}

	key := deriveKey(password, salt)

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

var rootCmd = &cobra.Command{
	Use:   "aegis",
	Short: "A simple file encryption and decryption tool",
	Long: ` Aegis is a CLI tool for encrypting and decrypting files based on the 
AES256 GCM (default) and CHACHA20 POLY1305 algorithms. It's based on 
minio/sio which implements the Data At Rest Encryption (DARE) format.`,
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt [input file] [output file]",
	Short: "Encrypt a file",
	Long: `Encrypt a file using a password.
If no output file is specified, it will use [input file].enc as the output file.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")
		inputFile := args[0]
		outputFile := inputFile + ".enc"
		if len(args) == 2 {
			outputFile = args[1]
		}
		err := encryptFile(inputFile, outputFile, password)
		if err != nil {
			fmt.Printf("Encryption error: %v\n", err)
			return
		}
		fmt.Printf("File encrypted successfully! Output: %s\n", outputFile)
	},
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt [input file] [output file]",
	Short: "Decrypt a file",
	Long: `Decrypt a file using a password.
If no output file is specified, it will use [input file].dec as the output file.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")
		inputFile := args[0]
		outputFile := inputFile + ".dec"
		if len(args) == 2 {
			outputFile = args[1]
		}
		err := decryptFile(inputFile, outputFile, password)
		if err != nil {
			fmt.Printf("Decryption error: %v\n", err)
			return
		}
		fmt.Printf("File decrypted successfully! Output: %s\n", outputFile)
	},
}

func main() {
	encryptCmd.Flags().StringP("password", "p", "", "Password for encryption (required)")
	encryptCmd.MarkFlagRequired("password")

	decryptCmd.Flags().StringP("password", "p", "", "Password for decryption (required)")
	decryptCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(encryptCmd, decryptCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
