package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	kitlog "github.com/go-kit/kit/log"

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

	// Logging
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stderr)

	// Create the redis client
	client, err := redis.NewClient(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Panic(err)
	}

	repo := redis.NewRepository(client)

	service := shortener.NewRedirectService(repo)

	// http server
	endpoints := api.MakeEndpoints(service, logger)
	ctx := context.Background()
	handler := api.NewHTTPServer(ctx, endpoints)

	// tcp server
	rpcHandler := api.NewRPC(service)

	err = rpc.Register(rpcHandler)
	if err != nil {
		log.Fatal("error registering API", err)
	}
	rpc.HandleHTTP()

	errs := make(chan error, 2)

	// http server
	go func() {
		fmt.Println("Running http server on port 8080")
		errs <- http.ListenAndServe(":8080", handler)
	}()

	go func() {
		fmt.Println("Running rpc server on port 8081")
		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			errs <- err
		}

		err = http.Serve(listener, nil)
		if err != nil {
			errs <- err
		}
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}
