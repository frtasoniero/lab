package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/frtasoniero/bidirectional-streaming-client/src/pb/shoppingcart"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("error on get connection. error: ", err)
	}
	defer conn.Close()

	client := shoppingcart.NewShoppingCartServiceClient(conn)
	stream, err := client.AddItem(context.Background())
	if err != nil {
		log.Fatalln("error on get channel to stream. error: ", err)
	}

	wait := make(chan struct{})
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				close(wait)
				return
			}
			if err != nil {
				log.Fatalln("error on recv. error: ", err)
			}

			fmt.Printf("<- response: %+v\n", response)
		}
	}()

	items := []shoppingcart.AddProduct{
		{ProductId: 1, Quantity: 2, PriceUnit: 5.0},
		{ProductId: 2, Quantity: 5, PriceUnit: 3.5},
		{ProductId: 3, Quantity: 1, PriceUnit: 2.5},
	}

	for _, value := range items {
		if err := stream.Send(&value); err != nil {
			log.Fatalln("error on send. error: ", err)
		}
		fmt.Printf("-> send: {productId:%+v, quantity:%+v, priceUnit:%+v}\n", value.ProductId, value.Quantity, value.PriceUnit)
		time.Sleep(1 * time.Second)
	}

	stream.CloseSend()

	<-wait
}
