// Модуль main представлен двумя подмодулями: Client и Server.
// В клиенте происходим демонстрация функционала нашего приложения.
package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"

	pb "gophkeeper/internal/app/proto"
	MyLogicClient "gophkeeper/internal/logic_client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	projectpath = "../../"
)

func main() {
	conn, err := grpc.Dial("localhost:3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("did not connect:", err)
	}
	defer conn.Close()

	c := pb.NewMyServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//Просто запрос на регистрацию.
	req1, err := c.Register(ctx, &pb.RegisterRequest{Login: "abcde@mail.ru", Password: "123456"})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req1\n", "Status:", req1.GetStatus(), "Cookie", req1.GetCookie())

	// Запрос на повторную регистрацию.
	req1ar, err := c.Register(ctx, &pb.RegisterRequest{Login: "abcde@mail.ru", Password: "123456"})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req1ar\n", "Status:", req1ar.GetStatus(), "Cookie", req1ar.GetCookie())

	//Просто запрос на регистрацию.
	req1ou, err := c.Register(ctx, &pb.RegisterRequest{Login: "123@mail.ru", Password: "qwertyuiop"})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req1ou\n", "Status:", req1ou.GetStatus(), "Cookie", req1ou.GetCookie())

	//Попытка аутентификации с правильными cookie.
	req2, err := c.Authentication(ctx, &pb.AuthenticationRequest{Login: "abcde@mail.ru", Password: "123456", Cookie: req1.GetCookie()})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req2\n", "Status:", req2.GetStatus(), "Cookie", req2.GetCookie())

	//Попытка аутентификации с неправильными cookie, подразумевается что можно зайти с неправильным куки, хендлер изменит его на правильный в случаи правильного пароля и логина.
	req2nc, err := c.Authentication(ctx, &pb.AuthenticationRequest{Login: "abcde@mail.ru", Password: "123456", Cookie: "asdfsdfsdfsdfsdfsdfsdfsdfsdfsdfsdfsdfsdf"})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req2nc\n", "Status:", req2nc.GetStatus(), "Cookie", req2nc.GetCookie())

	//Попытка аутентификации с неправильными логином, паролем и куки.
	req2tn, err := c.Authentication(ctx, &pb.AuthenticationRequest{Login: "123123@mail.ru", Password: "qwertyuiop", Cookie: "asdfsdfsdfsdfsdfsdfsdfsdfsdfsdfsdfsdfsdf"})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req2tn\n", "Status:", req2tn.GetStatus(), "Cookie", req2tn.GetCookie())

	//Смена мастер пароля.
	req3, err := c.EditMasterPass(ctx, &pb.EditMasterPassRequest{Cookie: req2.GetCookie(), Newpass: "654321"})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req3\n", "Status:", req3.GetStatus())

	//Смена мастер пароля.
	req3r, err := c.EditMasterPass(ctx, &pb.EditMasterPassRequest{Cookie: req2.GetCookie(), Newpass: "123456"})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req3r\n", "Status:", req3r.GetStatus())

	//Добавление данных в БД, данные текстового файла.
	buffFileTxtData := MyLogicClient.ReadFromFile(projectpath + "content/file.txt")
	strctFileTxtData := MyLogicClient.MyData{
		BytesData: hex.EncodeToString(buffFileTxtData),
		MetaInfo:  "File",
	}
	bFileTxtData, err := json.Marshal(strctFileTxtData)
	if err != nil {
		log.Println(err)
	}
	req4ft, err := c.AddData(ctx, &pb.AddDataRequest{Cookie: req2.GetCookie(), Data: bFileTxtData})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req4ft\n", "Status:", req4ft.GetStatus(), req4ft.GetId())

	//Добавление данных в БД, данные Логин/Пароль.
	logPass := MyLogicClient.LoginPassword{
		Login:    "abc@mail.ru",
		Password: "123456",
	}
	buffLogPass, err := json.Marshal(logPass)
	if err != nil {
		log.Println(err)
	}
	strctLogPass := MyLogicClient.MyData{
		BytesData: hex.EncodeToString(buffLogPass),
		MetaInfo:  "LoginPassword",
	}
	bLogPass, err := json.Marshal(strctLogPass)
	if err != nil {
		log.Println(err)
	}
	req4lp, err := c.AddData(ctx, &pb.AddDataRequest{Cookie: req2.GetCookie(), Data: bLogPass})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req4lp\n", "Status:", req4lp.GetStatus(), req4lp.GetId())

	//Добавление данных в БД, данные Банковской Карты.
	bankCard := MyLogicClient.BankCard{
		CardNumber: "0000 0000 0000 0000",
		Validity:   "31/99",
		CVV:        111,
		Pin:        0000,
	}
	buffBankCard, err := json.Marshal(bankCard)
	if err != nil {
		log.Println(err)
	}
	strctBankCard := MyLogicClient.MyData{
		BytesData: hex.EncodeToString(buffBankCard),
		MetaInfo:  "BankCard",
	}
	bBankCard, err := json.Marshal(strctBankCard)
	if err != nil {
		log.Println(err)
	}
	req4bc, err := c.AddData(ctx, &pb.AddDataRequest{Cookie: req2.GetCookie(), Data: bBankCard})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req4bc\n", "Status:", req4bc.GetStatus(), req4bc.GetId())

	//Обновление данных в БД, данные Логин/Пароль.
	newLogPass := MyLogicClient.LoginPassword{
		Login:    "new.abc@mail.ru",
		Password: "qwertyu",
	}
	bNewLogPass, err := json.Marshal(newLogPass)
	if err != nil {
		log.Println(err)
	}
	req5lp, err := c.UpdateData(ctx, &pb.UpdateDataRequest{Cookie: req2.GetCookie(), Data: bNewLogPass, Id: req4lp.GetId()})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req5lp\n", "Status:", req5lp.GetStatus())

	//Обновление данных в БД, данные Банковской карты.
	newBankCard := MyLogicClient.BankCard{
		CardNumber: "1111 1111 1111 1111",
		Validity:   "11/77",
		CVV:        999,
		Pin:        1111,
	}
	bNewBankCard, err := json.Marshal(newBankCard)
	if err != nil {
		log.Println(err)
	}
	req5bc, err := c.UpdateData(ctx, &pb.UpdateDataRequest{Cookie: req2.GetCookie(), Data: bNewBankCard, Id: req4bc.GetId()})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req5bc\n", "Status:", req5bc.GetStatus())

	//Запрос данных, данные логин/пароль.
	req6lp, err := c.GetData(ctx, &pb.GetDataRequest{Cookie: req2.GetCookie(), Id: req4lp.GetId()})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req6lp\n", "Status:", req6lp.GetStatus())
	MyLogicClient.CheckData(req6lp.GetData())

	//Запрос данных, данные банковской карты.
	req6bc, err := c.GetData(ctx, &pb.GetDataRequest{Cookie: req2.GetCookie(), Id: req4bc.GetId()})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req6bc\n", "Status:", req6bc.GetStatus())
	MyLogicClient.CheckData(req6bc.GetData())

	//Запрос данных, данные содержание стекстового файла.
	req6ft, err := c.GetData(ctx, &pb.GetDataRequest{Cookie: req2.GetCookie(), Id: req4ft.GetId()})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req6bc\n", "Status:", req6ft.GetStatus())
	MyLogicClient.CheckData(req6ft.GetData())

	//Удаление данных, данные содержание текстового файла.
	req7, err := c.Delete(ctx, &pb.DeleteRequest{Cookie: req2.GetCookie(), Id: req4ft.GetId()})
	if err != nil {
		log.Println("could not greet:", err)
	}
	log.Println("req7\n", "Status:", req7.GetStatus())
}
