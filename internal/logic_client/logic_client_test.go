package logic_client_test

import (
	"bytes"
	"testing"

	MyLogicClient "gophkeeper/internal/logic_client"
)

func TestReadFromFile(t *testing.T) {
	filePath := "../../content/file.txt"

	expectedContent := []byte("Hello World!")

	content := MyLogicClient.ReadFromFile(filePath)

	if !bytes.Equal(content, expectedContent) {
		t.Errorf("Expected content: %s, got: %s", expectedContent, content)
	}
}
