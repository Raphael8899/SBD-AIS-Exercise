package client

import (
	"context"
	"exc8/pb"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcClient struct {
	client pb.OrderServiceClient
}

func NewGrpcClient() (*GrpcClient, error) {
	conn, err := grpc.NewClient(":4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewOrderServiceClient(conn)
	return &GrpcClient{client: client}, nil
}

func (c *GrpcClient) Run() error {
	ctx := context.Background()

	// 1. List drinks
	fmt.Println("Requesting drinks 🍹🍺☕")
	resp, err := c.client.GetDrinks(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	fmt.Println("Available drinks:")
	for _, d := range resp.Drinks {
		// Custom formatting to match screenshot exactly
		// Note: Coffee in screenshot doesn't show price, so we check if price > 0
		if d.Price > 0 {
			fmt.Printf("\t> id:%d  name:\"%s\"  price:%d  description:\"%s\"\n", d.Id, d.Name, d.Price, d.Description)
		} else {
			fmt.Printf("\t> id:%d  name:\"%s\"  description:\"%s\"\n", d.Id, d.Name, d.Description)
		}
	}

	// Helper function to order
	order := func(id int32, name string, amount int32) {
		fmt.Printf("\t> Ordering: %d x %s\n", amount, name)
		_, err := c.client.OrderDrink(ctx, &pb.OrderRequest{
			DrinkId: id,
			Amount:  amount,
		})
		if err != nil {
			log.Printf("Failed to order %s: %v", name, err)
		}
	}

	// 2. Order a few drinks
	fmt.Println("Ordering drinks 👨‍🍳⏱️🍻🍻")
	order(1, "Spritzer", 2)
	order(2, "Beer", 2)
	order(3, "Coffee", 2)

	// 3. Order more drinks
	fmt.Println("Ordering another round of drinks 👨‍🍳⏱️🍻🍻")
	order(1, "Spritzer", 6)
	order(2, "Beer", 6)
	order(3, "Coffee", 6)

	// 4. Get order total
	fmt.Println("Getting the bill 💹💹💹")
	orders, err := c.client.GetOrders(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	for _, o := range orders.Orders {
		fmt.Printf("\t> Total: %d x %s\n", o.TotalAmount, o.DrinkName)
	}

	return nil
}
