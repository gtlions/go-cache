package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/gtlions/go-cache/cache"
	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	cache.UnimplementedCacheServer
}

var db map[string]string
var addr = flag.String("addr", "localhost:8080", "server addr")

func (s *server) Store(ctx context.Context, request *cache.StoreRequest) (*cache.StoreResponse, error) {
	log.Printf("Store:[ %v ]", request)
	db[request.Key] = request.Value
	return &cache.StoreResponse{Status: "OK"}, nil
}

func (s *server) Get(ctx context.Context, request *cache.GetRequest) (*cache.GetResponse, error) {
	log.Printf("Get:[ %v ]", request.Key)
	if v, ok := db[request.Key]; !ok {
		return &cache.GetResponse{}, fmt.Errorf("Key:[ %s ] not exists.", request.Key)
	} else {
		return &cache.GetResponse{Key: request.Key, Value: v}, nil
	}
}

func (s *server) StoreStream(stream cache.Cache_StoreStreamServer) error {
	i := 1
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Recv EOF\n")
			return nil
		} else if err != nil {
			log.Printf("Recv Err:[ %v ]\n", err)
			return err
		} else {
			log.Printf("Store:[ %v ]", request)
			db[request.Key] = request.Value
			if i%2 == 0 {
				if err := stream.Send(&cache.StoreResponse{Status: "OK"}); err != nil {
					log.Printf("Recv Err:[ %v ]\n", err)
					return err
				}
			}
		}
		i++
	}
}

func (s *server) List(ctx context.Context, request *emptypb.Empty) (*cache.ListResonse, error) {
	log.Printf("List")
	list := make([]*cache.GetResponse, 0)
	for k := range db {
		list = append(list, &cache.GetResponse{Key: k, Value: db[k]})
	}
	return &cache.ListResonse{List: list}, nil
}

func main() {
	flag.Parse()
	db = make(map[string]string)
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal("Listen err:", err)
	}
	ser := grpc.NewServer()
	cache.RegisterCacheServer(ser, &server{})
	log.Println("启动监听服务...", *addr)
	if err := ser.Serve(lis); err != nil {
		log.Fatal("Serve err:", err)
	}
}
