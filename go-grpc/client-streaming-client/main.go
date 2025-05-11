package main

import (
	"context"
	"fmt"
	"log"

	"github.com/frtasoniero/client-streaming-client/src/pb/calc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("error on new client. error: ", err)
	}
	defer conn.Close()

	client := calc.NewCalcServiceClient(conn)
	stream, err := client.Calc(context.Background())
	if err != nil {
		log.Fatalln("error on get channel to stream. error: ", err)
	}

	nums := []int32{1, 2, 3, 4}
	for _, value := range nums {
		if err := stream.Send(&calc.Input{Value: value}); err != nil {
			log.Fatalln("error on send. error: ", err)
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("error on close stream. error: ", err)
	}

	fmt.Printf("response: {\n  %+v\n}\n", response)
}
