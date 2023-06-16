package logic_server_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	MyLogicServer "gophkeeper/internal/logic_server"
)

func TestMakeCookie(t *testing.T) {
	login := "username"

	cookieHash, salt := MyLogicServer.MakeCookie(login)

	expectedSaltLength := 6
	expectedCookieHashLength := 32

	if len(salt) != expectedSaltLength {
		t.Errorf("Expected salt length: %d, got: %d", expectedSaltLength, len(salt))
	}

	if len(cookieHash) != expectedCookieHashLength {
		t.Errorf("Expected cookie hash length: %d, got: %d", expectedCookieHashLength, len(cookieHash))
	}
}

func TestRandSeq(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	expectedLength := 6

	salt := MyLogicServer.RandSeq(expectedLength)

	if len(salt) != expectedLength {
		t.Errorf("Expected salt length: %d, got: %d", expectedLength, len(salt))
	}

	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for _, char := range salt {
		if !strings.ContainsRune(letters, char) {
			t.Errorf("Invalid salt character: %c", char)
		}
	}
}
