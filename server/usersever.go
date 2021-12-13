package main

import (
	"context"
    "bytes"
	"log"
	pb "main/proto"
	"main/server/db"
	"net"
	"flag"
	"google.golang.org/grpc"
	"encoding/json"
    "io/ioutil"
    "net/http"
    "sort"
)

const (
	port = ":5050"
	HOST = "http://127.0.0.1:"
)

var userport, postport, followport *string

type UserServer struct {
	pb.UnimplementedUserServiceServer
}

type Userinfo struct {
	Password    string
	Profilename string
	Profileimg  string
}

type Post struct {
	Text string
	Time string
}

func (s *UserServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("Received: %v, %v", in.Username, in.Password)

	resp, err := http.Get(HOST+*userport+"/"+in.Username)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}

	var userinfo Userinfo
	json.Unmarshal(body, &userinfo)

	if userinfo != (Userinfo{}) {
		if userinfo.Password == in.Password {
			return &pb.LoginResponse{Username: in.Username, Profilename: userinfo.Profilename, Profileimg: userinfo.Profileimg, Status: "Success", Msg: ""}, nil
		} else {
			return &pb.LoginResponse{Status: "Fail", Msg: "password is not correct"}, nil
		}
	} else {
		return &pb.LoginResponse{Status: "Fail", Msg: "user does not exist"}, nil
	} 

}

func (s *UserServer) SignUp(ctx context.Context, in *pb.SignUpRequest) (*pb.CommResponse, error) {
	log.Printf("Received: %v, %v, %v, %v", in.Username, in.Password, in.Profilename, in.Profileimg)
	user := Userinfo {
		Password: in.Password,
		Profilename: in.Profilename,
		Profileimg: in.Profileimg,
	}
	userinfo, _ := json.Marshal(user)
	client := &http.Client{}
    req, err := http.NewRequest(http.MethodPut, HOST+*userport+"/"+in.Username, bytes.NewBuffer(userinfo))
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}
	if err := string(body); err == "" {
		return &pb.CommResponse{Status: "Success", Msg: err}, nil
	} else {
		log.Printf("Fail to insert the user: %v", err)
		return &pb.CommResponse{Status: "Fail", Msg: err}, nil
	}

}

func (s *UserServer) CreatePost(ctx context.Context, in *pb.PostRequest) (*pb.CommResponse, error) {
	post := Post {
		Text: in.Text,
		Time: in.Time,
	}
	newpost, _ := json.Marshal(post)
	client := &http.Client{}
    req, err := http.NewRequest(http.MethodPut, HOST+*postport+"/"+in.Username, bytes.NewBuffer(newpost))
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}
	if err := string(body); err == "" {
		return &pb.CommResponse{Status: "Success", Msg: err}, nil
	} else {
		log.Printf("Fail to create the post: %v", err)
		return &pb.CommResponse{Status: "Fail", Msg: err}, nil
	}

}

func (s *UserServer) GetPosts(ctx context.Context, in *pb.CommRequest) (*pb.PostResponse, error) {
	log.Printf("GetPosts Received: %v", in.Username)

	resp, err := http.Get(HOST+*postport+"/"+in.Username)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}

	var posts []Post
	json.Unmarshal(body, &posts)
	sort.Slice(posts, func(i, j int) bool {
	  return posts[i].Time > posts[j].Time
	})
	var post_responses []*pb.PostResponsePost
	for _, post := range posts {
		post_response := &pb.PostResponsePost{Username: in.Username, Profilename: "test", Profileimg: "", Text: post.Text, Img: "", Time: post.Time}
		post_responses = append(post_responses, post_response)		
	}
	log.Println("Query Post:", post_responses)
	return &pb.PostResponse{Posts: post_responses}, nil
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


func main() {
	userport = flag.String("userport", "10380", "user port")
	postport = flag.String("postport", "11380", "post port")
	followport = flag.String("followport", "12380", "follow port")
	flag.Parse()

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
