package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"sol.go/cwm/appfirebase"
	"sol.go/cwm/appgrpc"
	"sol.go/cwm/apphttp"
	"sol.go/cwm/appws"
	"sol.go/cwm/dao"
	"sol.go/cwm/pubsub"
	"strconv"
	"time"
)

func main() {
	// if we crash the go code, we get the file and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	grpcPort, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	listenGRPC := flag.Int("g", grpcPort, "wait for incoming connections")
	listenHTTP := flag.Int("h", httpPort, "wait for incoming connections")
	flag.Parse()

	grpcPort = *listenGRPC
	httpPort = *listenHTTP
	fmt.Println("grpcPort", grpcPort)
	fmt.Println("httpPort", httpPort)

	//connect to mongodb
	fmt.Println("Connecting to mongodb ...")
	db := dao.GetDataBase()
	mongoCtx, cancelMongo := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelMongo()
	// test mongo connection
	err = db.MongoClient.Ping(mongoCtx, nil)
	if err != nil {
		log.Fatalf("Failed to connect to mongodb: %v", err)
	}

	//connect to redis
	fmt.Println("Connecting to redis ...")
	cache := dao.GetCache()
	// test redis connection
	redisCtx, cancelRedis := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelRedis()
	_, err = cache.RedisClient.Ping(redisCtx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	listener, _, grpcServer, err := initServices(httpPort, grpcPort)
	if err != nil {
		log.Fatalf("Failed to init services: %v", err)
	}

	//Setup shutdown hook
	// Wait for Ctrl+C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block until a signal is received
	<-ch

	fmt.Println("Stopping grpc server")
	grpcServer.Stop()
	fmt.Println("Stopping grpc listener")
	_ = listener.Close()

	fmt.Println("Closing mongodb connection")
	_ = db.MongoClient.Disconnect(mongoCtx)
	fmt.Println("End of programm")
}

func initServices(httpPort int, grpcPort int) (net.Listener, *gin.Engine, *grpc.Server, error) {

	listener, grpcServer, err := appgrpc.StartServer(grpcPort)
	if err != nil {
		return nil, nil, nil, err
	}

	appFireBase := appfirebase.GetAppFireBase()
	appFireBase.Start()

	ws := appws.GetWS()
	ws.Start()

	router, err := apphttp.StartServer(httpPort, ws.Server)
	if err != nil {
		return nil, nil, nil, err
	}

	pubsub.StartSubscribe()

	return listener, router, grpcServer, nil
}
