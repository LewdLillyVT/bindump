package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func dumpBinaryToFiles(filePath string) error {
	// Get the current working directory (where the program is executed)
	execDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %v", err)
	}

	// Create output filenames
	outputTxtFileName := filepath.Base(filePath) + "_dump.txt"
	outputBinFileName := filepath.Base(filePath) + "_dump.bin"

	// Create full output paths
	outputTxtFilePath := filepath.Join(execDir, outputTxtFileName)
	outputBinFilePath := filepath.Join(execDir, outputBinFileName)

	// Create a text file to dump the binary data
	txtFile, err := os.Create(outputTxtFilePath)
	if err != nil {
		return fmt.Errorf("could not create output text file: %v", err)
	}
	defer txtFile.Close()

	// Create a binary file to dump the binary data
	binFile, err := os.Create(outputBinFilePath)
	if err != nil {
		return fmt.Errorf("could not create output binary file: %v", err)
	}
	defer binFile.Close()

	// Open the provided executable file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Read the file in chunks and write the binary data to both output files
	const chunkSize = 4096
	buf := make([]byte, chunkSize)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading file: %v", err)
		}
		if n == 0 {
			break
		}

		// Write the binary data in hexadecimal format to the text file
		for i := 0; i < n; i++ {
			_, err := fmt.Fprintf(txtFile, "%02X ", buf[i])
			if err != nil {
				return fmt.Errorf("error writing to output text file: %v", err)
			}
		}
		_, err = txtFile.WriteString("\n") // New line for readability
		if err != nil {
			return fmt.Errorf("error writing to output text file: %v", err)
		}

		// Write raw binary data to the binary file
		_, err = binFile.Write(buf[:n])
		if err != nil {
			return fmt.Errorf("error writing to output binary file: %v", err)
		}
	}

	fmt.Printf("Binary data dumped to:\n- %s\n- %s\n", outputTxtFilePath, outputBinFilePath)
	return nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Prompt user to drag and drop files into the console
	fmt.Println("Please drag and drop an executable file into this console, then press Enter to proceed:")
	filesInput, _ := reader.ReadString('\n')
	filesInput = strings.TrimSpace(filesInput) // Remove whitespace and newline characters
	files := filepath.SplitList(filesInput)

	// Trim quotes from file paths
	for i, filePath := range files {
		files[i] = strings.Trim(filePath, "\"")
	}

	if len(files) == 0 {
		fmt.Println("Error: No files provided. Please drag and drop an executable file.")
		return
	}

	// Iterate over each file path provided
	for _, filePath := range files {
		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Error: File %s does not exist. Skipping.\n", filePath)
			continue
		}

		// Dump binary data to files
		if err := dumpBinaryToFiles(filePath); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	// Keep the console open by prompting the user to press Enter to exit
	fmt.Println("\nDump complete. Press Enter to exit.")
	reader.ReadString('\n')
}
