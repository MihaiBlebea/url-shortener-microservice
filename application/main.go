package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/MihaiBlebea/url-shortener/api"
	"github.com/MihaiBlebea/url-shortener/repository/redis"
	"github.com/MihaiBlebea/url-shortener/shortener"
)

func main() {
	// client, err := mysql.NewClient(
	// 	os.Getenv("MYSQL_USER"),
	// 	os.Getenv("MYSQL_PASS"),
	// 	os.Getenv("MYSQL_HOST"),
	// 	os.Getenv("MYSQL_PORT"),
	// 	os.Getenv("MYSQL_DB"),
	// )
	// if err != nil {
	// 	log.Panic(err)
	// }
	// repo, err := mysql.NewRepository(client)
	// if err != nil {
	// 	log.Panic(err)
	// }

	client, err := redis.NewClient(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Panic(err)
	}

	repo := redis.NewRepository(client)

	service := shortener.NewRedirectService(repo)

	endpoints := api.MakeEndpoints(service)

	ctx := context.Background()
	handler := api.NewHTTPServer(ctx, endpoints)

	// http server
	go func() {
		http.ListenAndServe(":8080", handler)
	}()

	rpcHandler := api.NewRPC(service)

	// tcp server
	err = rpc.Register(rpcHandler)
	if err != nil {
		log.Fatal("error registering API", err)
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Listener error", err)
	}
	log.Printf("serving rpc on port %s", "8081")
	http.Serve(listener, nil)

	if err != nil {
		log.Fatal("error serving: ", err)
	}
}
