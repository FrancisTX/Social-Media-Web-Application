package main

import (
	"context"
    "bytes"
	"log"
	pb "main/proto"
	"net"
	"google.golang.org/grpc"
	"encoding/json"
    "io/ioutil"
    "net/http"
    "sort"
    "errors"
    "html/template"
)

const (
	port = ":5050"
	HOST = "http://127.0.0.1:"
	userport = "10380"
	postport = "11380"
	followport = "12380"
)


type UserServer struct {
	pb.UnimplementedUserServiceServer
}

type Userinfo struct {
	Password    string
	Profilename string
	Profileimg  template.URL
}

type Post struct {
	Text string
	Time string
	Img  template.URL
}

func (s *UserServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("Received: %v, %v", in.Username, in.Password)

	resp, err := http.Get(HOST+userport+"/"+in.Username)
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
    req, err := http.NewRequest(http.MethodPut, HOST+userport+"/"+in.Username, bytes.NewBuffer(userinfo))
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
		Img:  in.Img,
	}
	newpost, _ := json.Marshal(post)
	client := &http.Client{}
    req, err := http.NewRequest(http.MethodPut, HOST+postport+"/"+in.Username, bytes.NewBuffer(newpost))
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

	// Get all following users
	resp, _ := http.Get(HOST+followport+"/"+in.Username)
	body, _ := ioutil.ReadAll(resp.Body)
	var users []string
	json.Unmarshal(body, &users)
	users = append(users, in.Username)

	var post_responses []*pb.PostResponsePost

	for _, user := range users {
		resp_user, _ := http.Get(HOST+userport+"/"+user)
		body_user, _ := ioutil.ReadAll(resp_user.Body)
		var userinfo Userinfo
		json.Unmarshal(body_user, &userinfo)
		resp, _ := http.Get(HOST+postport+"/"+user)
		body, _ := ioutil.ReadAll(resp.Body)
		var posts []Post
		json.Unmarshal(body, &posts)
		for _, post := range posts {
			post_response := &pb.PostResponsePost{Username: user, Profilename: userinfo.Profilename, Profileimg: userinfo.Profileimg, Text: post.Text, Img: post.Img, Time: post.Time}
			post_responses = append(post_responses, post_response)
		}
	}
	sort.Slice(post_responses, func(i, j int) bool {
	  return post_responses[i].Time > post_responses[j].Time
	})
	return &pb.PostResponse{Posts: post_responses}, nil
}

func (s *UserServer) GetUserInfo(ctx context.Context, in *pb.CommRequest) (*pb.LoginResponse, error) {
	log.Printf("GetUserInfo Received: %v", in.Username)

	resp, err := http.Get(HOST+userport+"/"+in.Username)
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
		return &pb.LoginResponse{Username: in.Username, Profilename: userinfo.Profilename, Profileimg: userinfo.Profileimg}, nil
	} else {
		return &pb.LoginResponse{}, errors.New("user not found")
	}
}

func (s *UserServer) Follow(ctx context.Context, in *pb.FollowRequest) (*pb.CommResponse, error) {
	client := &http.Client{}
    req, err := http.NewRequest(http.MethodPut, HOST+followport+"/"+in.Username1, bytes.NewBuffer([]byte(in.Username2)))
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
		return &pb.CommResponse{Status: "Success", Msg: "Follow Finish"}, nil
	} else {
		return nil, errors.New(err)
	}
	
}

func (s *UserServer) Unfollow(ctx context.Context, in *pb.FollowRequest) (*pb.CommResponse, error) {
	client := &http.Client{}
    req, err := http.NewRequest(http.MethodDelete, HOST+followport+"/"+in.Username1, bytes.NewBuffer([]byte(in.Username2)))
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
		return &pb.CommResponse{Status: "Success", Msg: "UnFollow Finish"}, nil
	} else {
		return nil, errors.New(err)
	}
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
