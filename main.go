package main

import (
	"fmt"
	"github.com/Zhoangp/Cart-Service/config"
	"github.com/Zhoangp/Cart-Service/internal/delivery/http"
	"github.com/Zhoangp/Cart-Service/internal/repo"
	"github.com/Zhoangp/Cart-Service/internal/usecase"
	"github.com/Zhoangp/Cart-Service/pb"
	"github.com/Zhoangp/Cart-Service/pkg/database/mysql"
	"github.com/Zhoangp/Cart-Service/pkg/utils"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	env := os.Getenv("ENV")
	fileName := "config/config-local.yml"
	if env == "app" {
		fileName = "config/config-app.yml"
	}
	cf, err := config.LoadConfig(fileName)
	if err != nil {
		log.Fatalln("Failed at config", err)
	}
	db, err := mysql.NewMysql(cf)
	if err != nil {
		fmt.Println("Connection to database failed")
	}
	//client, err := redis2.NewRedis(cf)
	//if err != nil {
	//	log.Fatal("Cannot connect to redis", err)
	//}
	//ctx := context.Background()
	//
	//err = client.Set(ctx, "foo", "bar", 0).Err()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//val, err := client.Get(ctx, "foo").Result()
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println("foo", val)
	//}

	hasher := utils.NewHasher("courses", 3)
	cartRepo := repo.NewCartRepo(db)
	cartUsecase := usecase.NewCartUseCase(cf, cartRepo, hasher)
	cartHandle := http.NewCartHandler(cartUsecase, cf)

	lis, err := net.Listen("tcp", ":"+cf.Service.Port)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Cart Svc on", cf.Service.Port)
	grpcServer := grpc.NewServer()
	pb.RegisterCartServiceServer(grpcServer, cartHandle)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}

}
