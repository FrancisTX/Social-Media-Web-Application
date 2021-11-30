package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "main/proto"
)

const (
	address = "localhost:5050"
)

func Login(args map[string]string) (map[string]string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Login

	r, err := c.Login(ctx, &pb.LoginRequest{Username: args["username"], Password: args["password"]})
	if r.Status == "Success" {
		return map[string]string{"username":r.Username, "profilename":r.Profilename, "profileimg":r.Profileimg, "status": r.Status}
	} else {
		return map[string]string{"status": r.Status, "msg": r.Msg}
	}
	
}

func SignUp(args map[string]string) (map[string]string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// SignUp

	r, err := c.SignUp(ctx, &pb.SignUpRequest{Username: args["username"], Password: args["password"], Profilename: args["profilename"], Profileimg: args["profileimg"]})
	return map[string]string{"status": r.Status, "msg": r.Msg}
}