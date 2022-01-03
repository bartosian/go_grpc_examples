package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could'n connect %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// doUnary(c)

	// doServerSreaming(c)

	doClientStreaming(c)
}

func doServerSreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Server Streaming RPC")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Kiryl",
			LastName:  "Bartashevich",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("Print reponse of GreetManyTimes %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Client Streaming RPC")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.GreetRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Kiryl",
				},
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.GreetRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Michael",
				},
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.GreetRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Ann",
				},
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.GreetRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Petr",
				},
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.GreetRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Bob",
				},
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.GreetRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Yuri",
				},
			},
		},
	}

	stream, err := c.LonGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreat: %v", err)
	}

	for _, req := range requests {
		stream.Send(req)
		time.Sleep(4e9)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from LongGreet %v", err)
	}

	fmt.Printf("Received resppnse from LonGreet %v", res)
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
