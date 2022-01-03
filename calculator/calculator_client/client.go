package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can't create client connection %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(c)
	// doServerSreaming(c)
	doClientStreaming(c)
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Sending ClientStreaming Request.")

	numbers := []calculatorpb.ComputeAverageRequest{
		calculatorpb.ComputeAverageRequest{Number: 2},
		calculatorpb.ComputeAverageRequest{Number: 6},
		calculatorpb.ComputeAverageRequest{Number: 12},
		calculatorpb.ComputeAverageRequest{Number: 9},
		calculatorpb.ComputeAverageRequest{Number: 7},
		calculatorpb.ComputeAverageRequest{Number: 29},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error sending ComputeAverage Request %v", err)
	}

	for _, number := range numbers {
		stream.Send(&number)
		time.Sleep(1e9)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response from server %v", err)
	}

	fmt.Printf("Received response from ComputeAverage: %v", res)

}

func doServerSreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Sending CalculatorService Streaming Request.")
	req := calculatorpb.PrimeNumberRequest{
		Number: 80,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error calling Streaming Api %v", err)
	}

	for {
		num, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		fmt.Println("Received value from stream %v", num)
	}
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Sending Calcu;atorService Request.")
	req := calculatorpb.CalculatorRequest{
		Request: &calculatorpb.Sum{
			Num_1: 23,
			Num_2: 56,
		},
	}

	res, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("Failed to send Sum request %v", err)
	}

	fmt.Printf("Got response from server %v", res.GetResponse())
}
