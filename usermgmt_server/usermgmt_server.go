package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "github.com/mel-ak/unary/usermgmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	port=":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
	}
}
type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
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
	var users_list *pb.UserList =&pb.UserList{}
	var user_id int32 = int32(rand.Intn(1000))
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
	
	readBytes, err:= ioutil.ReadFile("users.json")
	if err != nil {
		if os.IsNotExist(err){
			log.Println("File Not found: ", err.Error())
			log.Println("Creating New File")
			users_list.Users= append(users_list.Users, created_user)
			jsonBytes, err := protojson.Marshal(users_list)
			if err!= nil {
                log.Println("Error marshaling json", err)
            }
			if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
				log.Println("Error writing to file", err)
			}

			return created_user,nil
		}
		log.Println("Error reading from file", err)
    }

    if err := protojson.Unmarshal(readBytes, users_list); err != nil {
		log.Println("Error unmarshalling json", err)
	}

	users_list.Users= append(users_list.Users, created_user)
	jsonBytes, err := protojson.Marshal(users_list)
	if err!= nil {
        log.Println("Error marshaling json", err)
    }
	if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err!= nil {
        log.Println("Error writing to file", err)
    }
	return created_user,nil
}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	jsonBytes, err := ioutil.ReadFile("users.json")
	if err!= nil {
		log.Println("Error reading from file", err)
        return nil, err
	}

	var users_list *pb.UserList = &pb.UserList{}
	if err := protojson.Unmarshal(jsonBytes, users_list); err!= nil {
        log.Println("Error unmarshalling json", err)
        return nil, err
    }
	return users_list,nil
}

func main() {
	// Instantiate
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("Failed to initialize user management server: %v", err)
	}
}