package main

import (
	"context"
	"flag"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"grpc-gateway-rnd/internal/core/util"
	pb_user_v1 "grpc-gateway-rnd/internal/gen/proto/user/v1"
	"log/slog"
	"net"
	"os/signal"
	"sync"
)

var (
	addr = flag.String("addr", "localhost:9000", "grpc addr to listen")
)

type UserServiceGrpc struct {
	pb_user_v1.UnimplementedUserServer

	mu sync.Mutex
	db map[string]*pb_user_v1.UserModel
}

func NewUserServiceGrpc() *UserServiceGrpc {
	return &UserServiceGrpc{db: make(map[string]*pb_user_v1.UserModel)}
}

func (s *UserServiceGrpc) Create(_ context.Context, req *pb_user_v1.UserCreateRequest) (*pb_user_v1.UserCreateResponse, error) {
	m := pb_user_v1.UserModel_builder{
		Id:   util.Ref(uuid.New().String()),
		Name: util.Ref(util.Or(req.GetName(), "Chuck Norris")),
		Age:  util.Ref(util.Or(req.GetAge(), 99999)),
	}.Build()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.db[m.GetId()] = m

	return pb_user_v1.UserCreateResponse_builder{User: m}.Build(), nil
}

func (s *UserServiceGrpc) Get(_ context.Context, req *pb_user_v1.UserGetRequest) (*pb_user_v1.UserGetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	m := s.db[req.GetId()]
	return pb_user_v1.UserGetResponse_builder{User: m}.Build(), nil
}

func (s *UserServiceGrpc) Delete(_ context.Context, req *pb_user_v1.UserDeleteRequest) (*pb_user_v1.UserDeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	m := s.db[req.GetId()]
	delete(s.db, req.GetId())

	return pb_user_v1.UserDeleteResponse_builder{User: m}.Build(), nil
}

func main() {
	ctx, cancelFunc := signal.NotifyContext(context.Background())
	defer cancelFunc()
	wg := new(sync.WaitGroup)

	flag.Parse()

	l, err := new(net.ListenConfig).Listen(ctx, "tcp", *addr)
	util.UnreachableErrWithMessage(err, "failed to listen GRPC addr")

	grpcServer := grpc.NewServer()
	pb_user_v1.RegisterUserServer(grpcServer, NewUserServiceGrpc())

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Info("serving GRPC...", "addr", *addr)
		err := grpcServer.Serve(l)
		util.UnreachableErrWithMessage(err, "failed to serve GRPC server")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		grpcServer.GracefulStop()
		grpcServer.Stop()
	}()

	<-ctx.Done()
	wg.Wait()
}
