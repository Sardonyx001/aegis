package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

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
		err := aegis.encryptFile(inputFile, outputFile, password)
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
