package main

import (
	"context"
	"log"
	"time"

	"github.com/sundayyogurt/product_service/internal/api/grpc/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := product.NewProductServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetProductById(ctx, &product.GetProductRequest{Id: 1})
	if err != nil {
		panic(err)
	}

	log.Printf("Response: %+v\n", resp)
}
