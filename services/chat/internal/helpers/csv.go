package helpers

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func GenerateCSV(fileName string, numAccounts int, publicFolder string) error {
	filePath := filepath.Join(publicFolder, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"first_name", "last_name", "email"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for i := 1; i <= numAccounts; i++ {
		firstName := fmt.Sprintf("Test_%d", i)
		lastName := fmt.Sprintf("Test_%d", i)
		email := fmt.Sprintf("test_email_%d@tests.com", i)
		record := []string{firstName, lastName, email}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	fmt.Printf("File created: %s accounts: %d\n", fileName, numAccounts)
	return nil
}

func CountCSVRows(file *multipart.FileHeader) (int, error) {
	uploadedFile, err := file.Open()
	headerCount := 1
	if err != nil {
		return 0, err
	}
	defer uploadedFile.Close()

	reader := csv.NewReader(bufio.NewReader(uploadedFile))

	var rowCount int
	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return rowCount, err
		}
		rowCount++
	}

	return rowCount - headerCount, nil
}
