package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/frtasoniero/server-streaming-client/src/pb/department"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("error on get connection. error: ", err)
	}
	defer conn.Close()

	client := department.NewDepartmentServiceClient(conn)
	stream, err := client.ListPerson(context.Background(), &department.ListPersonRequest{DepartmentId: 2})
	if err != nil {
		log.Fatalln("error on get channel to stream. error: ", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("error on recv. error: ", err)
		}

		fmt.Printf("response: %+v\n", response)
	}
}
