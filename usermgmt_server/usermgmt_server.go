package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/mel-ak/unary/usermgmt"
	"google.golang.org/grpc"
)

const (
	port=":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		user_list: &pb.UserList{},
	}
}
type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
	user_list *pb.UserList
}

func (server *UserManagementServer) Run() error {
	lis,err := net.Listen("tcp",port)
	if err != nil {
		log.Fatalf("Failed to listen on %v ", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s,server)

	log.Printf("Server listening on %v", lis.Addr())
	return s.Serve(lis)
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Println("Recieved ", in.GetName())
	var user_id int32 = int32(rand.Intn(1000))
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
	s.user_list.Users = append(s.user_list.Users, created_user)
	return created_user,nil
}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	return s.user_list,nil
}

func main() {
	// Instantiate
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("Failed to initialize user management server: %v", err)
	}
}