package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "homework-1/pkg/api/v1"
	"log"
)

// Test GRPC client
func main() {
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewAdminServiceClient(conn)

	ctx := context.Background()
	productId := uint64(1)
	pageNum := uint64(0)
	pageSize := uint64(4)
	productListRequest := pb.ProductListRequest{Page: &pageNum, Size: &pageSize}

	createProductRequest := pb.ProductCreateRequest{
		Name:     "first product",
		Price:    23,
		Quantity: 7,
	}
	if response, err := client.ProductCreate(ctx, &createProductRequest); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("create: %v \n", response)
		productId = response.GetId()
	}

	if response, err := client.ProductList(ctx, &productListRequest); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("list: %v \n", response)
	}

	if response, err := client.ProductGet(ctx, &pb.ProductGetRequest{Id: productId}); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("get: %v \n", response)
	}

	createProductRequest = pb.ProductCreateRequest{
		Name:     "second product",
		Price:    100,
		Quantity: 10,
	}
	if response, err := client.ProductCreate(ctx, &createProductRequest); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("create: %v \n", response)
		productId = response.GetId()
	}

	if response, err := client.ProductList(ctx, &productListRequest); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("list: %v \n", response)
	}

	updateProductRequest := pb.ProductUpdateRequest{
		Id:       productId,
		Name:     "New name",
		Price:    20,
		Quantity: 30,
	}
	if response, err := client.ProductUpdate(ctx, &updateProductRequest); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("update: %v \n", response)
	}

	if response, err := client.ProductList(ctx, &productListRequest); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("list: %v \n", response)
	}

	if response, err := client.ProductDelete(ctx, &pb.ProductDeleteRequest{Id: productId}); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("delete: %v \n", response)
	}

	if response, err := client.ProductList(ctx, &productListRequest); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("list: %v \n", response)
	}
}
