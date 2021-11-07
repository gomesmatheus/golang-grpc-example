package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/codeedu/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}
	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j.@j.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "j.@j.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}
		fmt.Println("Status:", stream.Status, " - ", stream.GetUser())

	}

	fmt.Println(responseStream)

}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "m1",
			Name:  "Matheus 1",
			Email: "matheus1@gmail.com",
		},
		&pb.User{
			Id:    "m2",
			Name:  "Matheus 2",
			Email: "matheus2@gmail.com",
		},
		&pb.User{
			Id:    "m3",
			Name:  "Matheus 3",
			Email: "matheus3@gmail.com",
		},
		&pb.User{
			Id:    "m4",
			Name:  "Matheus 4",
			Email: "matheus4@gmail.com",
		},
		&pb.User{
			Id:    "m5",
			Name:  "Matheus 5",
			Email: "matheus5@gmail.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "m1",
			Name:  "Matheus 1",
			Email: "matheus1@gmail.com",
		},
		&pb.User{
			Id:    "m2",
			Name:  "Matheus 2",
			Email: "matheus2@gmail.com",
		},
		&pb.User{
			Id:    "m3",
			Name:  "Matheus 3",
			Email: "matheus3@gmail.com",
		},
		&pb.User{
			Id:    "m4",
			Name:  "Matheus 4",
			Email: "matheus4@gmail.com",
		},
		&pb.User{
			Id:    "m5",
			Name:  "Matheus 5",
			Email: "matheus5@gmail.com",
		},
	}

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	wait := make(chan int)

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Receiving User: %v com status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait

}
