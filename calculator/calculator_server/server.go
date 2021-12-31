package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Getting CalculatorRequest: %v", req)
	num_1, num_2 := req.GetRequest().GetNum_1(), req.GetRequest().GetNum_2()
	sum := num_1 + num_2
	res := calculatorpb.CalculatorResponse{
		Response: sum,
	}

	return &res, nil
}

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Getting primeNumberDecompositionRequest: %v", req)
	k := int32(2)
	N := req.GetNumber()

	for N > 1 {
		if N%k == 0 {
			res := calculatorpb.PrimeNumberResponse{
				Response: k,
			}
			stream.Send(&res)
			N = N / k
		} else {
			k += int32(1)
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Can't create listener %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
