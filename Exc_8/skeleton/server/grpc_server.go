package server

import (
	"context"
	"exc8/pb"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GRPCService struct {
	pb.UnimplementedOrderServiceServer
	// In-memory storage
	drinks    []*pb.Drink
	orders    map[int32]int32 // Maps DrinkID to TotalAmount
	orderLock sync.Mutex
}

func StartGrpcServer() error {
	srv := grpc.NewServer()

	grpcService := &GRPCService{
		orders: make(map[int32]int32),
		drinks: []*pb.Drink{
			{Id: 1, Name: "Spritzer", Price: 2, Description: "Wine with soda"},
			{Id: 2, Name: "Beer", Price: 3, Description: "Hagenberger Gold"},
			{Id: 3, Name: "Coffee", Description: "Mifare isn't that secure"},
		},
	}

	pb.RegisterOrderServiceServer(srv, grpcService)

	lis, err := net.Listen("tcp", ":4000")
	if err != nil {
		return err
	}

	err = srv.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func (s *GRPCService) GetDrinks(ctx context.Context, _ *emptypb.Empty) (*pb.DrinkList, error) {
	return &pb.DrinkList{Drinks: s.drinks}, nil
}

func (s *GRPCService) OrderDrink(ctx context.Context, req *pb.OrderRequest) (*wrapperspb.BoolValue, error) {
	s.orderLock.Lock()
	defer s.orderLock.Unlock()

	found := false
	for _, d := range s.drinks {
		if d.Id == req.DrinkId {
			found = true
			break
		}
	}

	if found {
		s.orders[req.DrinkId] += req.Amount
		return &wrapperspb.BoolValue{Value: true}, nil
	}

	return &wrapperspb.BoolValue{Value: false}, nil
}

func (s *GRPCService) GetOrders(ctx context.Context, _ *emptypb.Empty) (*pb.OrderList, error) {
	s.orderLock.Lock()
	defer s.orderLock.Unlock()

	var orderItems []*pb.OrderItem

	for _, d := range s.drinks {
		if amount, ok := s.orders[d.Id]; ok && amount > 0 {
			orderItems = append(orderItems, &pb.OrderItem{
				DrinkName:   d.Name,
				TotalAmount: amount,
			})
		}
	}

	return &pb.OrderList{Orders: orderItems}, nil
}
