package order

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderServiceGRPCClient struct {
	client OrderServiceClient
}

func NewOrderServiceGRPCClient(orderSvcAddress string) (*OrderServiceGRPCClient, error) {
	// เปิด connection ไปหา order service
	client, err := grpc.NewClient(orderSvcAddress, grpc.WithTransportCredentials(insecure.NewCredentials())) // "localhost:50051"
	if err != nil {
		return nil, err
	}

	oderServiceClient := NewOrderServiceClient(client)

	return &OrderServiceGRPCClient{oderServiceClient}, nil
}

func (o *OrderServiceGRPCClient) GetOrderIserOrders(userid int64) (*GetOrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel() // ถ้าเกิน 5 วิ หรือเสร็จแล้ว cancel ทันที

	res, err := o.client.GetOrders(ctx, &GetOrderRequest{
		UserId: userid,
	}) // เรียก gRPC ไปที่ order service
	if err != nil {
		return nil, err
	}

	return res, nil
}
