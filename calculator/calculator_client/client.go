package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"

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
	doServerSreaming(c)
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
