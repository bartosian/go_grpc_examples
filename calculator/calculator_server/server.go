package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Getting Average: %v", stream)
	sum := 0
	len := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			response := calculatorpb.ComputeAverageResponse{
				Average: float64(sum / len),
			}

			return stream.SendAndClose(&response)
		}

		if err != nil {
			log.Fatalf("Error processing ComputeAverage request: %v", err)
		}

		sum += int(req.GetNumber())
		len++
	}
}

func (s *server) CalculateMaximum(stream calculatorpb.CalculatorService_CalculateMaximumServer) error {
	fmt.Printf("Calculate Maximum request: %v", stream)
	max := int32(0)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error getting request from client: %v\n", err)
			return err
		}

		number := req.GetNumber()
		if number > max {
			max = number
		}

		errSend := stream.Send(&calculatorpb.CalculateMaximumResponse{
			Maximum: max,
		})

		if errSend != nil {
			log.Fatalf("Error sending response: %v\n", err)
			return err
		}
	}
}

func (s *server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("Calculate Aquare Root")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number %v", number),
		)
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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
