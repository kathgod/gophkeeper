// Пакет logic_server реализует вспомогательные функции и типы данных для работы сервера.
package logic_server

import (
	"crypto/sha256"
	mr "math/rand"
)

// MailPass Структура для хранения логина и пароля.
type MailPass struct {
	Mail string `json:"mail"`
	Pass string `json:"pass"`
}

// CookieSalt Структура для хранения куки и соли.
type CookieSalt struct {
	Cookie string `json:"cookie"`
	Salt   string `json:"salt,omitempty"`
}

// MakeCookie Функция создание новой куки.
func MakeCookie(login string) ([32]byte, string) {
	salt := RandSeq(6)
	var s = append([]byte(login), []byte(salt)...)

	return sha256.Sum256(s), salt
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandSeq для создания случайной соли наперед заданной длинны.
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mr.Intn(len(letters))]
	}
	return string(b)
}

// IDCookie Структура для сохраненных данных.
type IDCookie struct {
	ID     int    `json:"id"`
	Cookie string `json:"cookie"`
}
