package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// doClientStreaming(c)

	// doBiDiStreaming(c)

	doUnaryWithDeadline(c)
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

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting BiDi Streaming RPC")

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Kiryl",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Michael",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ann",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Olga",
			},
		},
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling GreatEveryone: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1e9)
		}

		stream.CloseSend()
	}()

	go func() {
		defer close(waitc)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}

			fmt.Printf("received: %v", res.GetResult())
		}
	}()

	<-waitc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Unary With Deadline RPC")
	req := greetpb.GreetGreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Kiryl",
			LastName:  "Bartashevich",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, &req)
	if err != nil {
		statusErr, ok := status.FromError(err)

		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit!")
			} else {
				fmt.Printf("Unexpceted error: %v", statusErr)
			}
		} else {
			log.Fatalf("Error while calling GreetWithDeadline RPC: %v", err)
		}

		log.Fatalf("Error while calling Greet RPC: %v", err)

		return
	}

	log.Printf("Response from Greet: %v", res.Result)
}
