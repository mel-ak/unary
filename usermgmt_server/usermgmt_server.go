package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v4"
	pb "github.com/mel-ak/unary/usermgmt"
	"google.golang.org/grpc"
)

const (
	port=":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
	}
}
type UserManagementServer struct {
	conn *pgx.Conn
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

func (server *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	createSql := `
	create table if not exists users(
		id SERIAL PRIMARY KEY,
		name TEXT,
		age int
	);
	`
	_, err := server.conn.Exec(context.Background(), createSql)
	if err!= nil {
        log.Printf("Failed to create users table: %v", err)
		os.Exit(1)
        return nil, err
    }

	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge()}

	tx, err := server.conn.Begin(context.Background())
	if err!= nil {
        log.Printf("Failed to start transaction: %v", err)
        return nil, err
    }
	_, err = tx.Exec(context.Background(), "insert into users(name,age) values ($1,$2)", created_user.Name, created_user.Age)
	if err!= nil {
        log.Printf("Failed to insert user: %v", err)
        tx.Rollback(context.Background())
        return nil, err
    }
	tx.Commit(context.Background())




	return created_user,nil
}

func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	var users_list *pb.UserList = &pb.UserList{}
	rows, err := server.conn.Query(context.Background(), "select * from users")
	if err != nil {
		log.Printf("Failed to get users: %v", err)
        return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err!= nil {
            log.Printf("Failed to scan row: %v", err)
            return nil, err
        }
		users_list.Users = append(users_list.Users, &user)
	}
	
	return users_list,nil
}

func main() {
	// Connect to Postgres
	conn, err := pgx.Connect(context.Background(), "host=localhost port=5431 user=user password=password dbname=db sslmode=disable")
    if err!= nil {
        log.Printf("Failed to connect to database: %v", err)
        os.Exit(1)
    }
	log.Println("Database connected")
    defer conn.Close(context.Background())
	// Instantiate
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	user_mgmt_server.conn = conn
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("Failed to initialize user management server: %v", err)
	}
}