package main

import (
	"context"
	pb "github.com/deniz/grpc-demo/proto"
	"io"
	"log"
	"strconv"
	"time"
)

func saveUser(client pb.GreetServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.SaveUser(ctx, &pb.SaveRequest{
		Name:    "Deniz",
		SurName: "Kaplan",
		Email:   "deniz@kaplan.com",
	})
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("%s", res.Id)
}

func GetUsersEmail(client pb.GreetServiceClient) {
	log.Printf("Client Streaming started")
	stream, err := client.GetUsersEmail(context.Background())
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	var list []pb.SaveRequest

	list = append(list, pb.SaveRequest{Name: "admin", Email: "deniz@kaplan.com", SurName: "admin"})
	list = append(list, pb.SaveRequest{Name: "admin1", Email: "deniz1@kaplan.com", SurName: "admin1"})
	list = append(list, pb.SaveRequest{Name: "admin2", Email: "deniz2@kaplan.com", SurName: "admin2"})

	for _, item := range list {
		if err := stream.Send(&item); err != nil {
			log.Fatalf("Error while sending %v", err)
		}
		//log.Printf("Sent request with name: %s", item.Name)
		//time.Sleep(2 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	log.Printf("Client Streaming finished")
	if err != nil {
		log.Fatalf("Error while receiving %v", err)
	}
	log.Printf("email %v surname %v", res.Email, res.SurName)
}

func GetIdUserMail(client pb.GreetServiceClient) {
	log.Printf("Streaming started")
	user := &pb.User{
		Name:    "deniz",
		SurName: "kaplan",
		Email:   "deniz@kaplan.com",
		Id:      1,
	}
	stream, err := client.GetIdUserMail(context.Background(), user)
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming %v", err)
		}
		log.Println(message)
	}

	log.Printf("Streaming finished")
}

func GetAlllUser(client pb.GreetServiceClient) {
	log.Printf("Bidirectional Streaming started")
	stream, err := client.GetAlllUser(context.Background())
	if err != nil {
		log.Fatalf("Could not send names: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while streaming %v", err)
			}
			log.Println(message)
		}
		close(waitc)
	}()

	for i := 0; i < 5; i++ {
		req := &pb.SaveRequest{
			Name:    "admin" + strconv.Itoa(i),
			SurName: "admin" + strconv.Itoa(i),
			Email:   "deniz" + strconv.Itoa(i),
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending %v", err)
		}
	}

	stream.CloseSend()
	<-waitc
	log.Printf("Bidirectional Streaming finished")

}
