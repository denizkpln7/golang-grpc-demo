package main

import (
	"context"
	"github.com/deniz/grpc-demo/proto"
	"io"
	"log"
	"strconv"
)

func (s *userServer) SaveUser(ctx context.Context, request *proto.SaveRequest) (*proto.User, error) {
	println("hi")
	return &proto.User{
		Id:      10,
		Email:   request.Email,
		SurName: request.SurName,
		Name:    request.Name,
	}, nil
}

func (s *userServer) GetUsersEmail(stream proto.GreetService_GetUsersEmailServer) error {
	var user proto.User
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&user)
		}
		if err != nil {
			return err
		}
		log.Printf("Got request with name : %v", req.Name)

		if req.Email == "deniz@deniz1@kaplan.com.com" {
			user = proto.User{
				Id:      10,
				Email:   req.Email,
				SurName: req.SurName,
				Name:    req.Name,
			}
		}
	}
}

func (s *userServer) GetIdUserMail(req *proto.User, stream proto.GreetService_GetIdUserMailServer) error {
	log.Printf("Got request with user : %v", req.Name)

	for i := 0; i < 5; i++ {
		res := &proto.SaveRequest{
			Email:   req.Email,
			Name:    req.Name + strconv.Itoa(i),
			SurName: req.SurName + strconv.Itoa(i),
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (s *userServer) GetAlllUser(stream proto.GreetService_GetAlllUserServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Got request with name : %v", req.Name)
		res := &proto.User{
			Id:      10,
			Email:   req.Email,
			SurName: req.SurName,
			Name:    req.Name,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
}
