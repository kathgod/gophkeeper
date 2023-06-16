// Модуль main представлен двумя подмодулями: Client и Server.
// В сервере происходим хост функционала нашего приложения.
package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	MyHandler "gophkeeper/internal/app"
	pb "gophkeeper/internal/app/proto"
)

func main() {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterMyServiceServer(s, &MyHandler.UserServer{})
	log.Printf("server listening at %v", listen.Addr())
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
