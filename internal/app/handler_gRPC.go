// Пакет handler представляет собой набор функции которые могут вызываться из клиента, и которые реализуют различный функционал нашего сервиса.
package handler

import (
	"context"
	"encoding/hex"
	pb "gophkeeper/internal/app/proto"
	"log"

	MyStorage "gophkeeper/internal/database"
	MyLogicServer "gophkeeper/internal/logic_server"
)

// UserServer для имплементации моих хендлеров на сервер.
type UserServer struct {
	pb.UnimplementedMyServiceServer
}

// Дуть до корня проекта.
const (
	projectpath = "../../"
)

// Register хендлер для регистрации пользователя.
func (s *UserServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	clientdb := projectpath + "internal/database/client.db"
	serverdb := projectpath + "internal/database/server.db"

	db := MyStorage.OpenDB(clientdb)
	defer MyStorage.Sync(clientdb, serverdb)
	defer db.Close()

	mailPass := MyLogicServer.MailPass{
		Mail: req.Login,
		Pass: req.Password,
	}

	cookieValue, salt := MyLogicServer.MakeCookie(req.GetPassword())

	cookieSalt := MyLogicServer.CookieSalt{
		Cookie: hex.EncodeToString(cookieValue[:]),
		Salt:   salt,
	}

	status := 200

	err, buff := MyStorage.RegisterDB(db, mailPass, cookieSalt)

	if err != nil {
		log.Println(err)
		cookieSalt.Cookie = buff
		status = 409
	}

	return &pb.RegisterResponse{Status: int64(status), Cookie: cookieSalt.Cookie}, nil
}

// Authentication хендлер аутентификации.
func (s *UserServer) Authentication(ctx context.Context, req *pb.AuthenticationRequest) (*pb.AuthenticationResponse, error) {
	db := MyStorage.OpenDB(projectpath + "/internal/database/client.db")
	defer db.Close()

	mailPassGet := MyLogicServer.MailPass{
		Mail: req.GetLogin(),
		Pass: req.GetPassword(),
	}

	cookieSaltGet := MyLogicServer.CookieSalt{
		Cookie: req.GetCookie(),
	}

	ckSltDB, err := MyStorage.GetCookieDB(db, mailPassGet)
	if err != nil {
		return &pb.AuthenticationResponse{Status: 404, Cookie: "NOT FOUND"}, nil
	}

	if cookieSaltGet.Cookie == ckSltDB.Cookie {
		return &pb.AuthenticationResponse{Status: 200, Cookie: cookieSaltGet.Cookie}, nil
	} else {
		return &pb.AuthenticationResponse{Status: 200, Cookie: ckSltDB.Cookie}, nil
	}
}

// EditMasterPass хендлер смены мастер пароля.
func (s *UserServer) EditMasterPass(ctx context.Context, req *pb.EditMasterPassRequest) (*pb.EditMasterPassResponse, error) {
	clientdb := projectpath + "internal/database/client.db"
	serverdb := projectpath + "internal/database/server.db"

	db := MyStorage.OpenDB(clientdb)
	defer MyStorage.Sync(clientdb, serverdb)
	defer db.Close()

	authStatus := MyStorage.Authorization(db, req.GetCookie())
	if !authStatus {
		return &pb.EditMasterPassResponse{Status: 401}, nil
	}

	MyStorage.EditMasterPass(db, req.GetCookie(), req.GetNewpass())

	return &pb.EditMasterPassResponse{Status: 200}, nil
}

// AddData хендлер добавления данных в систему.
func (s *UserServer) AddData(ctx context.Context, req *pb.AddDataRequest) (*pb.AddDataResponse, error) {
	clientdb := projectpath + "internal/database/client.db"
	serverdb := projectpath + "internal/database/server.db"

	db := MyStorage.OpenDB(clientdb)
	defer MyStorage.Sync(clientdb, serverdb)
	defer db.Close()

	authStatus := MyStorage.Authorization(db, req.GetCookie())
	if !authStatus {
		return &pb.AddDataResponse{Status: 401, Id: -1}, nil
	}

	ID := MyStorage.AddData(db, req.GetCookie(), req.GetData())

	return &pb.AddDataResponse{Status: 200, Id: int64(ID)}, nil
}

// UpdateData хендлер обновления данных в системе.
func (s *UserServer) UpdateData(ctx context.Context, req *pb.UpdateDataRequest) (*pb.UpdateDataResponse, error) {
	clientdb := projectpath + "internal/database/client.db"
	serverdb := projectpath + "internal/database/server.db"

	db := MyStorage.OpenDB(clientdb)
	defer MyStorage.Sync(clientdb, serverdb)
	defer db.Close()

	authStatus := MyStorage.Authorization(db, req.GetCookie())
	if !authStatus {
		return &pb.UpdateDataResponse{Status: 401}, nil
	}

	MyStorage.UpdateData(db, req.GetCookie(), req.GetData(), int(req.GetId()))

	return &pb.UpdateDataResponse{Status: 200}, nil
}

// GetData хендлер запроса данных.
func (s *UserServer) GetData(ctx context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	db := MyStorage.OpenDB(projectpath + "/internal/database/client.db")
	defer db.Close()

	authStatus := MyStorage.Authorization(db, req.GetCookie())
	if !authStatus {
		return &pb.GetDataResponse{Status: 401}, nil
	}

	bDataInfo, err := MyStorage.GetData(db, req.GetCookie(), int(req.GetId()))
	if err != nil {
		return &pb.GetDataResponse{Status: 401}, nil
	}

	return &pb.GetDataResponse{Status: 200, Data: bDataInfo}, nil
}

// Delete хендлер удаления данных.
func (s *UserServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	clientdb := projectpath + "internal/database/client.db"
	serverdb := projectpath + "internal/database/server.db"

	db := MyStorage.OpenDB(clientdb)
	defer db.Close()

	authStatus := MyStorage.Authorization(db, req.GetCookie())
	if !authStatus {
		return &pb.DeleteResponse{Status: 401}, nil
	}

	err := MyStorage.Delete(db, req.GetCookie(), int(req.GetId()))
	if err != nil {
		return &pb.DeleteResponse{Status: 502}, nil
	}

	sdb := MyStorage.OpenDB(serverdb)
	defer sdb.Close()

	err = MyStorage.Delete(sdb, req.GetCookie(), int(req.GetId()))
	if err != nil {
		log.Println(err)
	}

	return &pb.DeleteResponse{Status: 200}, nil
}
