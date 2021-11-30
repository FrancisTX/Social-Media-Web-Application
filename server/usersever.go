package main

import (
	"context"
	"log"
	"net"
	"fmt"
	"main/server/auth"
	"main/server/db"
	"google.golang.org/grpc"
	pb "main/proto"
)

const (
	port = ":5050"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("Received: %v, %v", in.Username, in.Password)
	if user, err := auth.Auth(in.Username, in.Password); err == "" {
		return &pb.LoginResponse{Username: user.Username, Profilename: user.ProfileName, Profileimg: user.ProfileImg, Status: "Success", Msg: err}, nil
	} else {
		return &pb.LoginResponse{Status: "Fail", Msg: err}, nil
	}
	
}

func (s *UserServer) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.CommResponse, error) {
	log.Printf("Received: %v, %v, %v, %v", in.Username, in.Password, in.Profilename, in.Profileimg)
	if err := db.InsertUser(in.Username, in.Password, in.Profilename, in.Profileimg); err == "" {
		return &pb.CommResponse{Status: "Success", Msg: err}, nil
	} else {
		return &pb.CommResponse{Status: "Fail", Msg: err}, nil
	}
	
}

func (s *UserServer) CreatePost(ctx context.Context, in *pb.PostRequest) (*pb.CommResponse, error) {
	if err := db.CreatePost(in.Username, in.Profilename, in.Profileimg, in.Text, in.Img, in.Time); err == "" {
		return &pb.CommResponse{Status: "Success", Msg: err}, nil
	} else {
		return &pb.CommResponse{Status: "Fail", Msg: err}, nil
	}
	
}

func (s *UserServer) GetPosts(ctx context.Context, in *pb.CommRequest) (*pb.PostResponse, error) {
	log.Printf("Received: %v", in.Username)
	rows, err := db.QueryPost(in.Username)
	var posts []*pb.PostResponsePost
	for rows.Next() {
		post := new(pb.PostResponsePost)
		err := rows.Scan(&post.Id, &post.Username, &post.Profilename, &post.Profileimg, &post.Text, &post.Img, &post.Time)
		if err != nil {
	        fmt.Println("Error while query posts: %v", err)
	    }
	    posts = append(posts, post)
	}
	return &pb.PostResponse{Posts: posts}, err
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
