package main

import (
	_ "github.com/golang/protobuf/proto"
	pb "github.com/jackgardner/go-ledger/proto"
	"golang.org/x/net/context"
	"strconv"
)

type LedgerServer struct {
	Transactions map[string]*pb.Transaction
}

func (s *LedgerServer) CreateTransaction(ctx context.Context, request *pb.CreateTransactionRequest) (*pb.Transaction, error) {
	var transactionId = strconv.Itoa(len(s.Transactions) + 1)
	var newTransaction = &pb.Transaction{Success: true, TransactionId: transactionId, AmountInPence: request.AmountInPence}
	s.Transactions[transactionId] = newTransaction
	return newTransaction, nil
}

func (s *LedgerServer) GetTransaction(ctx context.Context, request *pb.GetTransactionRequest) (*pb.Transaction, error) {
	return s.Transactions[request.TransactionId], nil
}

func (s *LedgerServer) GetTransactions(ctx context.Context, request *pb.ListTransactionsRequest) (*pb.TransactionsReply, error) {
	values := make([]*pb.Transaction, len(s.Transactions))
	i := 0
	for _, k := range s.Transactions {
		values[i] = k
		i++
	}

	return &pb.TransactionsReply{
		Transactions: values,
		LedgerId:     request.SourceLedgerId,
	}, nil
}
