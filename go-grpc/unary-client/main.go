package main

import (
	"context"
	"fmt"
	"log"

	"github.com/frtasoniero/unary-client/src/pb/products"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5005", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("error on get client. error: ", err)
	}
	defer conn.Close()

	findAllProducts(conn)
	//createProduct(conn)
}

func findAllProducts(conn *grpc.ClientConn) {
	productClient := products.NewProductServiceClient(conn)
	productList, err := productClient.FindAll(context.Background(), &products.Product{})
	if err != nil {
		log.Fatalln("error on list products. error: ", err)
	}
	fmt.Printf("products: %+v\n", productList)
}

func createProduct(conn *grpc.ClientConn) {
	newProduct := &products.Product{
		Name:        "Keyboard",
		Description: "HyperX Alloy Origins",
		Price:       450,
		Quantity:    2,
	}

	productClient := products.NewProductServiceClient(conn)
	newProduct, err := productClient.Create(context.Background(), newProduct)
	if err != nil {
		log.Fatalln("error on create product. error: ", err)
	}

	fmt.Printf("Product created!\n%+v\n", newProduct)
}
