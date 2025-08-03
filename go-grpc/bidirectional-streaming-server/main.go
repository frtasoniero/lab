package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/frtasoniero/bidirectional-streaming-server/src/pb/shoppingcart"
	"google.golang.org/grpc"
)

type server struct {
	shoppingcart.ShoppingCartServiceServer
}

func (s *server) AddItem(srv shoppingcart.ShoppingCartService_AddItemServer) error {
	var quantityItems int32 = 0
	var totalPrice float64 = 0.0

	for {
		newItem, err := srv.Recv()
		if err == io.EOF {
			return srv.Send(&shoppingcart.ShoppingCartTotal{
				QuantityItems: quantityItems,
				TotalPrice:    totalPrice,
			})
		}
		if err != nil {
			return fmt.Errorf("error on recv. error: %+v", err)
		}

		quantityItems += newItem.GetQuantity()
		totalPrice += float64(newItem.GetPriceUnit() * float64(newItem.GetQuantity()))

		if err := srv.Send(&shoppingcart.ShoppingCartTotal{
			QuantityItems: quantityItems,
			TotalPrice:    totalPrice,
		}); err != nil {
			return fmt.Errorf("error on send. error: %v", err)
		}
	}
}

func main() {
	fmt.Println("Starting grpc server...")

	listener, err := net.Listen("tcp", ":5005")
	if err != nil {
		log.Fatalln("error on get listener. error: ", err)
	}

	s := grpc.NewServer()
	shoppingcart.RegisterShoppingCartServiceServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalln("error on serve. error: ", err)
	}
}
