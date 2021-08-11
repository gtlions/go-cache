package main

import (
	"context"
	"flag"
	"io"
	"log"
	"strconv"

	"github.com/gtlions/go-cache/cache"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var addr = flag.String("addr", "localhost:8080", "server addr")

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Dial err:", err)
	}
	defer conn.Close()
	c := cache.NewCacheClient(conn)

	// Store
	for i := 1; i <= 10; i++ {
		k := "key" + strconv.Itoa(i)
		v := "value" + strconv.Itoa(i)
		storeResp, err := c.Store(context.Background(), &cache.StoreRequest{Key: k, Value: v})
		if err != nil {
			log.Fatal("Store err:", err)
		}
		log.Printf("Store Key:[ %s ], %v\n", k, storeResp)
	}
	// StoreStream
	storeStream, err := c.StoreStream(context.Background())
	for i := 1; i <= 10; i++ {
		k := "keyStream" + strconv.Itoa(i)
		v := "valueStream" + strconv.Itoa(i)
		log.Printf("Send Key:[ %s ]\n", k)
		err := storeStream.Send(&cache.StoreRequest{Key: k, Value: v})
		if err != nil {
			log.Fatal("storeStream send err:", err)
		}
	}
	storeStream.CloseSend()
	for {
		resp, err := storeStream.Recv()
		if err == io.EOF {
			log.Printf("storeStream Recv EOF\n")
			break
		} else if err != nil {
			log.Fatal("storeStream recv err:", err)
		} else {
			log.Printf("Store Key %v\n", resp)
		}
	}

	// Get
	getKey := "name"
	getResp, err := c.Get(context.Background(), &cache.GetRequest{Key: getKey})
	if err != nil {
		log.Printf("Get err:[ %s ]\n", err.Error())
	} else {
		log.Printf("Get Key:\n[ %s ], %v\n", getKey, getResp)
	}

	// List
	listResp, err := c.List(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatal("List err:", err)
	}
	log.Printf("List Key:\n %v\n", listResp)
}
