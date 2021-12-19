package client

import (
	"context"
	"log"
	"time"
	"html/template"

	pb "main/proto"

	"google.golang.org/grpc"
)

const (
	address = "localhost:5050"
)

func BuildConnections() (*grpc.ClientConn, pb.UserServiceClient, context.Context, context.CancelFunc) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100 * 1024 * 1024), grpc.MaxCallSendMsgSize(100 * 1024 * 1024)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	return conn, c, ctx, cancel
}

func Login(args map[string]string) (map[string]string, template.URL) {
	conn, c, ctx, cancel := BuildConnections()
	defer conn.Close()
	defer cancel()

	// Login

	r, _ := c.Login(ctx, &pb.LoginRequest{Username: args["username"], Password: args["password"]})
	if r.Status == "Success" {
		return map[string]string{"username": r.Username, "profilename": r.Profilename, "status": r.Status}, r.Profileimg
	} else {
		return map[string]string{"status": r.Status, "msg": r.Msg}, template.URL("")
	}

}

func SignUp(args map[string]string, profileimg template.URL) map[string]string {
	conn, c, ctx, cancel := BuildConnections()
	defer conn.Close()
	defer cancel()

	// SignUp

	r, _ := c.SignUp(ctx, &pb.SignUpRequest{Username: args["username"], Password: args["password"], Profilename: args["profilename"], Profileimg: profileimg})
	return map[string]string{"status": r.Status, "msg": r.Msg}
}

func CreatePost(args map[string]string, img template.URL) map[string]string {
	conn, c, ctx, cancel := BuildConnections()
	defer conn.Close()
	defer cancel()

	// CreatePost

	r, _ := c.CreatePost(ctx, &pb.PostRequest{Username: args["username"], Text: args["text"], Img: img, Time: args["time"]})
	return map[string]string{"status": r.Status, "msg": r.Msg}
}

func GetPosts(args map[string]string) ([]*pb.PostResponsePost, error) {
	conn, c, ctx, cancel := BuildConnections()
	defer conn.Close()
	defer cancel()

	// GetPost
	r, err := c.GetPosts(ctx, &pb.CommRequest{Username: args["username"]})
	if err != nil {
		return nil, err
	}
	return r.Posts, nil
}

func GetUserInfo(username string) (*pb.LoginResponse, error) {
	conn, c, ctx, cancel := BuildConnections()
	defer conn.Close()
	defer cancel()

	//Get User Info

	r, err := c.GetUserInfo(ctx, &pb.CommRequest{Username: username})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func Follow(username1, username2 string) (*pb.CommResponse, error) {
	conn, c, ctx, cancel := BuildConnections()
	defer conn.Close()
	defer cancel()

	//Follow
	r, err := c.Follow(ctx, &pb.FollowRequest{Username1: username1, Username2: username2})
	if err != nil {
		return nil, err
	}
	return r, nil
}

func Unfollow(username1, username2 string) (*pb.CommResponse, error) {
	conn, c, ctx, cancel := BuildConnections()
	defer conn.Close()
	defer cancel()

	r, err := c.Unfollow(ctx, &pb.FollowRequest{Username1: username1, Username2: username2})
	if err != nil {
		return nil, err
	}
	return r, nil
}
