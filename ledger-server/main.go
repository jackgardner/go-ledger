package main

import (
	"github.com/namsral/flag"
	"fmt"
	_ "github.com/golang/protobuf/proto"
	pb "github.com/jackgardner/go-ledger/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
	"io/ioutil"
	"encoding/json"
)
var (
	port int
)

func newServer() *LedgerServer {
	var transactions map[string]*pb.Transaction

	file, err := ioutil.ReadFile("test_data.json")
	if err != nil {
		grpclog.Fatalf("Failed to load default ledger data: %v", err)
	}

	if err := json.Unmarshal(file, &transactions); err != nil {
		grpclog.Fatalf("Failed to load default ledgar data: %v", err)
	}

	return &LedgerServer{
		Transactions: transactions,
	}
}

func main() {

	flag.IntVar(&port, "port", 3000, "Service port")
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}
	grpclog.Printf("Listening on port: %d", port)

	var opts []grpc.ServerOption

	// TODO Server cert

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterLedgerServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
