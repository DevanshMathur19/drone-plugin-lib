package harness

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	CIErrorMessageKey = "CI_ERROR_MESSAGE"
	CIErrorCodeKey    = "CI_ERROR_CODE"
	CIMetadataFileEnv = "CI_ERROR_METADATA"
)

// SetError sets the error message and error code and writes them to the CI_ERROR_METADATA file
func SetError(message, code string) error {
	if err := WriteEnvToMetadataFile(CIErrorMessageKey, message); err != nil {
		return err
	}
	return WriteEnvToMetadataFile(CIErrorCodeKey, code)
}

// WriteEnvToMetadataFile writes a key-value pair to the CI_ERROR_METADATA file
func WriteEnvToMetadataFile(key, value string) error {
	metadataFilePath := os.Getenv(CIMetadataFileEnv)
	if metadataFilePath == "" {
		return fmt.Errorf("environment variable %s is not set", CIMetadataFileEnv)
	}

	// Check the extension of the metadata file (.env or .out)
	ext := strings.ToLower(filepath.Ext(metadataFilePath))

	var content string
	if ext == ".env" {
		// Write in .env format (KEY=VALUE)
		content = fmt.Sprintf("%s=%s\n", key, value)
	} else if ext == ".out" {
		// Write in .out format (export KEY="VALUE")
		content = fmt.Sprintf("%s \"%s\"\n", key, value)
	} else {
		return fmt.Errorf("unsupported file extension: %s", ext)
	}

	return writeToFile(metadataFilePath, content)
}

// Helper function to append content to the file
func writeToFile(filename, content string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
