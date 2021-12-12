package main

import (
	"context"
	"log"
	"strings"
	pb "main/proto"
	"main/server/auth"
	"main/server/db"
	"main/server/storage"
	"net"
	"flag"
	"time"
	"google.golang.org/grpc"
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
		log.Printf("Fail to insert the user: %v", err)
		return &pb.CommResponse{Status: "Fail", Msg: err}, nil
	}

}

func (s *UserServer) CreatePost(ctx context.Context, in *pb.PostRequest) (*pb.CommResponse, error) {
	if err := db.CreatePost(in.Username, in.Text, in.Img, in.Time); err == "" {
		return &pb.CommResponse{Status: "Success", Msg: err}, nil
	} else {
		log.Printf("Fail to create the post: %v", err)
		return &pb.CommResponse{Status: "Fail", Msg: err}, nil
	}

}

func (s *UserServer) GetPosts(ctx context.Context, in *pb.CommRequest) (*pb.PostResponse, error) {
	log.Printf("GetPosts Received: %v", in.Username)
	rows, err := db.QueryPost(in.Username)
	if err != nil {
		return nil, err
	}
	var posts []*pb.PostResponsePost
	for _, row := range rows {
		for row.Next() {
			post := new(pb.PostResponsePost)
			err := row.Scan(&post.Username, &post.Profilename, &post.Profileimg, &post.Text, &post.Img, &post.Time)
			if err != nil {
				log.Println("Error while query posts:", err)
				continue
			}
			posts = append(posts, post)
		}
	}
	log.Println("Query Post:", posts)
	return &pb.PostResponse{Posts: posts}, nil
}

func (s *UserServer) GetUserInfo(ctx context.Context, in *pb.CommRequest) (*pb.LoginResponse, error) {
	log.Printf("GetUserInfo Received: %v", in.Username)
	var err error
	user, err := db.QueryUser(in.Username)
	if err != nil {
		log.Println("GetUserInfo QueryUser Fault:", err)
		return nil, err
	}
	log.Println("GetUserInfo query user:", user)
	return &pb.LoginResponse{Username: user.Username, Profilename: user.ProfileName, Profileimg: user.ProfileImg}, nil
}

func (s *UserServer) Follow(ctx context.Context, in *pb.FollowRequest) (*pb.CommResponse, error) {
	err := db.Follow(in.Username1, in.Username2)
	if err != nil {
		return nil, err
	}
	return &pb.CommResponse{Status: "Success", Msg: "Follow Finish"}, nil
}

func (s *UserServer) Unfollow(ctx context.Context, in *pb.FollowRequest) (*pb.CommResponse, error) {
	err := db.Unfollow(in.Username1, in.Username2)
	if err != nil {
		return nil, err
	}
	return &pb.CommResponse{Status: "Success", Msg: "Unfollow Finish"}, nil
}



var userkvs *storage.Userkvstore
var postkvs *storage.Postkvstore
var followkvs *storage.Followkvstore

func main() {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	port := flag.String("port", "9121", "server port")
	flag.Parse()

	proposeUserC := make(chan string)
	//proposePostC := make(chan string)
	//proposeFollowC := make(chan string)
	defer close(proposeUserC)
	//defer close(proposePostC)
	//defer close(proposeFollowC)

	getSnapshotUser := func() ([]byte, error) { return userkvs.GetSnapshot() }
	commitUserC, errorUserC, snapshotterReadyUser := storage.NewRaftNode(*id, strings.Split(*cluster, ","), getSnapshotUser, proposeUserC)
	//commitPostC, errorPostC, snapshotterReady := storage.NewRaftNode(*id, strings.Split(*cluster, ","), getSnapshot, proposePostC)
	//commitFollowC, errorFollowC, snapshotterReady := storage.NewRaftNode(*id, strings.Split(*cluster, ","), getSnapshot, proposeFollowC)

	userkvs = storage.NewUserKVStore(<-snapshotterReadyUser, proposeUserC, commitUserC, errorUserC)
	//postkvs = storage.NewKVStore(<-snapshotterReady, proposePostC, commitPostC, errorPostC)
	//followkvs = storage.NewKVStore(<-snapshotterReady, proposeFollowC, commitFollowC, errorFollowC)

	userkvs.Propose("test", storage.Userinfo{Password: "test", Profilename: "bot", Profileimg: ""})
	time.Sleep(5 * time.Second)
	if v, ok := userkvs.Lookup("test"); ok {
		log.Printf(v.Profilename)
	} else {
		log.Printf("***********************************************")
	}
	
	lis, err := net.Listen("tcp", ":" + *port)

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
