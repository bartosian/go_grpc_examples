package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could'n connect %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Unary RPC")
	req := greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Kiryl",
			LastName:  "Bartashevich",
		},
	}
	res, err := c.Greet(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}
