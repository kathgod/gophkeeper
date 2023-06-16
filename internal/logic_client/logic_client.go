// Пакет logic_client реализует вспомогательные функции и типы данных для работы клиента.
package logic_client

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// MyData Структура хранения данных пользователя с мето-информацией.
type MyData struct {
	BytesData string `json:"bytesdata"`
	MetaInfo  string `json:"metainfo,omitempty"`
}

// LoginPassword Структура для работы с логином и паролем.
type LoginPassword struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// BankCard Структура для работы данными банковской карты.
type BankCard struct {
	CardNumber string `json:"cardnumber"`
	Validity   string `json:"validity"`
	CVV        int    `json:"cvv"`
	Pin        int    `json:"pin"`
}

// ReadFromFile Функция чтения из файла.
func ReadFromFile(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println(err)
	}

	return content
}

// CheckData Функция проверки что вернул хендлер запроса данных.
func CheckData(reqData []byte) {
	var DataInfoR6 MyData

	if err := json.Unmarshal(reqData, &DataInfoR6); err != nil {
		log.Println(err)
	}

	if DataInfoR6.MetaInfo == "LoginPassword" {
		var Buff LoginPassword
		b, err := hex.DecodeString(DataInfoR6.BytesData)
		if err != nil {
			log.Println(err)
		}
		if err := json.Unmarshal(b, &Buff); err != nil {
			log.Println(err)
		}
		log.Println(Buff.Login, Buff.Password)
	} else if DataInfoR6.MetaInfo == "BankCard" {
		var Buff BankCard
		b, err := hex.DecodeString(DataInfoR6.BytesData)
		if err != nil {
			log.Println(err)
		}
		if err := json.Unmarshal(b, &Buff); err != nil {
			log.Println(err)
		}
		log.Println(Buff.CardNumber, Buff.Validity, Buff.CVV, Buff.Pin)
	} else {
		var Buff string
		b, err := hex.DecodeString(DataInfoR6.BytesData)
		if err != nil {
			log.Println(err)
		}
		Buff = string(b)
		log.Println(Buff)
	}
}
